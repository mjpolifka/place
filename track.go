//go:build windows

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type openWindow struct {
	HWND        uintptr `json:"hwnd"`
	ProcessName string  `json:"process_name"`
}

type trackedWindow struct {
	HWND        uint64 `json:"hwnd"`
	ProcessName string `json:"process_name"`
	CreatedSeq  uint64 `json:"created_seq"`
}

type trackerState struct {
	NextCreatedSeq uint64                   `json:"next_created_seq"`
	Windows        map[string]trackedWindow `json:"windows"`
}

func runTracker() error {
	fmt.Println("Starting window tracker...")
	statePath, err := stateFilePath()
	if err != nil {
		return err
	}

	state, err := loadTrackerState(statePath)
	if err != nil {
		return err
	}

	if err := refreshTrackerState(state); err != nil {
		return err
	}
	if err := saveTrackerState(statePath, state); err != nil {
		return err
	}
	fmt.Printf("Tracking %d windows. State file: %s\n", len(state.Windows), statePath)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	defer signal.Stop(sig)

	for {
		select {
		case <-ticker.C:
			if err := refreshTrackerState(state); err != nil {
				fmt.Println("tracker refresh failed:", err)
				continue
			}
			if err := saveTrackerState(statePath, state); err != nil {
				fmt.Println("tracker save failed:", err)
				continue
			}
		case <-sig:
			fmt.Println("Tracker stopping...")
			return nil
		}
	}
}

func stateFilePath() (string, error) {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		var err error
		base, err = os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("could not determine state file location: %w", err)
		}
	}
	return filepath.Join(base, "place", "state.json"), nil
}

func loadTrackerState(path string) (*trackerState, error) {
	state := &trackerState{NextCreatedSeq: 1, Windows: map[string]trackedWindow{}}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return state, nil
		}
		return nil, fmt.Errorf("failed reading state file %s: %w", path, err)
	}

	if len(data) == 0 {
		return state, nil
	}

	if err := json.Unmarshal(data, state); err != nil {
		return nil, fmt.Errorf("failed parsing state file %s: %w", path, err)
	}
	if state.Windows == nil {
		state.Windows = map[string]trackedWindow{}
	}
	if state.NextCreatedSeq == 0 {
		state.NextCreatedSeq = 1
	}
	return state, nil
}

func saveTrackerState(path string, state *trackerState) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed creating state directory: %w", err)
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed serializing tracker state: %w", err)
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return fmt.Errorf("failed writing temp state file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed replacing state file: %w", err)
	}
	return nil
}

func refreshTrackerState(state *trackerState) error {
	windows, err := enumerateOpenWindows()
	if err != nil {
		return err
	}

	seen := make(map[string]struct{}, len(windows))
	for _, w := range windows {
		key := hwndKey(w.HWND)
		seen[key] = struct{}{}
		if existing, ok := state.Windows[key]; ok {
			existing.ProcessName = w.ProcessName
			state.Windows[key] = existing
			continue
		}

		state.Windows[key] = trackedWindow{
			HWND:        uint64(w.HWND),
			ProcessName: w.ProcessName,
			CreatedSeq:  state.NextCreatedSeq,
		}
		state.NextCreatedSeq++
	}

	for key := range state.Windows {
		if _, ok := seen[key]; !ok {
			delete(state.Windows, key)
		}
	}

	return nil
}

func hwndKey(hwnd uintptr) string {
	return strconv.FormatUint(uint64(hwnd), 10)
}

func sortHWNDsByTrackedOrder(hwnds []uintptr) []uintptr {
	statePath, err := stateFilePath()
	if err != nil {
		return hwnds
	}
	state, err := loadTrackerState(statePath)
	if err != nil {
		return hwnds
	}

	type row struct {
		hwnd uintptr
		seq  uint64
		idx  int
	}

	rows := make([]row, 0, len(hwnds))
	for i, hwnd := range hwnds {
		seq := uint64(^uint64(0))
		if tw, ok := state.Windows[hwndKey(hwnd)]; ok {
			seq = tw.CreatedSeq
		}
		rows = append(rows, row{hwnd: hwnd, seq: seq, idx: i})
	}

	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].seq == rows[j].seq {
			return rows[i].idx < rows[j].idx
		}
		return rows[i].seq < rows[j].seq
	})

	result := make([]uintptr, 0, len(rows))
	for _, r := range rows {
		result = append(result, r.hwnd)
	}
	return result
}

func enumerateOpenWindows() ([]openWindow, error) {
	user32DLL := syscall.NewLazyDLL("user32.dll")
	kernel32DLL := syscall.NewLazyDLL("kernel32.dll")

	enumWindowsProc := user32DLL.NewProc("EnumWindows")
	getWindowThreadProcessIDProc := user32DLL.NewProc("GetWindowThreadProcessId")
	isWindowVisibleProc := user32DLL.NewProc("IsWindowVisible")

	openProcessProc := kernel32DLL.NewProc("OpenProcess")
	queryFullProcessImageNameWProc := kernel32DLL.NewProc("QueryFullProcessImageNameW")
	closeHandleProc := kernel32DLL.NewProc("CloseHandle")

	const processQueryLimitedInformation = 0x1000
	var windows []openWindow
	var callbackErr error

	callback := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
		visible, _, _ := isWindowVisibleProc.Call(hwnd)
		if visible == 0 {
			return 1 // continue
		}

		var pid uint32
		getWindowThreadProcessIDProc.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
		if pid == 0 {
			return 1 // continue
		}

		processHandle, _, _ := openProcessProc.Call(processQueryLimitedInformation, 0, uintptr(pid))
		if processHandle == 0 {
			return 1 // continue
		}
		defer closeHandleProc.Call(processHandle)

		buf := make([]uint16, syscall.MAX_PATH)
		size := uint32(len(buf))
		ret, _, _ := queryFullProcessImageNameWProc.Call(
			processHandle,
			0,
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&size)),
		)
		if ret == 0 {
			return 1 // continue
		}

		exePath := syscall.UTF16ToString(buf[:size])
		exeName := strings.ToLower(filepath.Base(exePath))
		windows = append(windows, openWindow{HWND: hwnd, ProcessName: exeName})
		return 1 // continue
	})

	ret, _, err := enumWindowsProc.Call(callback, 0)
	if ret == 0 && len(windows) == 0 && err != syscall.Errno(0) {
		callbackErr = err
	}
	if callbackErr != nil {
		return nil, fmt.Errorf("failed to enumerate windows: %w", callbackErr)
	}

	return windows, nil
}
