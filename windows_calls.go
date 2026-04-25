//go:build windows

package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

type displayDimensions struct {
	Height        int
	Width         int
	DisplayNumber int
	IsMain        bool
	Left          int
	Top           int
	Right         int
	Bottom        int
}

type rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type monitorInfo struct {
	CbSize    uint32
	RcMonitor rect
	RcWork    rect
	DwFlags   uint32
}

func moveWindow(windowName string, x int, y int, width int, height int) error {
	user32Dll := syscall.NewLazyDLL("user32.dll")
	moveWindowProc := user32Dll.NewProc("MoveWindow")
	setForegroundWindowProc := user32Dll.NewProc("SetForegroundWindow")
	bringWindowToTopProc := user32Dll.NewProc("BringWindowToTop")
	showWindowProc := user32Dll.NewProc("ShowWindow")

	hwndList, err := getHWNDs(windowName)
	if err != nil {
		return err
	}
	repaint := 1

	ret, _, err := showWindowProc.Call(hwndList[0], 9)
	if ret == 0 {
		return err
	}

	ret, _, err = moveWindowProc.Call(
		hwndList[0],
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(repaint),
	)
	if ret == 0 {
		return err
	}

	bringWindowToTopProc.Call(hwndList[0])
	setForegroundWindowProc.Call(hwndList[0])

	return nil
}

func getHWNDs(targetProcessName string) ([]uintptr, error) {
	user32DLL := syscall.NewLazyDLL("user32.dll")
	kernel32DLL := syscall.NewLazyDLL("kernel32.dll")

	enumWindowsProc := user32DLL.NewProc("EnumWindows")
	getWindowThreadProcessIDProc := user32DLL.NewProc("GetWindowThreadProcessId")
	isWindowVisibleProc := user32DLL.NewProc("IsWindowVisible")

	openProcessProc := kernel32DLL.NewProc("OpenProcess")
	queryFullProcessImageNameWProc := kernel32DLL.NewProc("QueryFullProcessImageNameW")
	closeHandleProc := kernel32DLL.NewProc("CloseHandle")

	const processQueryLimitedInformation = 0x1000
	var foundHWNDs []uintptr
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

		processHandle, _, _ := openProcessProc.Call(
			processQueryLimitedInformation,
			0,
			uintptr(pid),
		)
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
		if exeName == targetProcessName {
			foundHWNDs = append(foundHWNDs, hwnd)
		}

		return 1 // continue
	})

	ret, _, err := enumWindowsProc.Call(callback, 0)
	if ret == 0 && len(foundHWNDs) == 0 {
		if err != syscall.Errno(0) {
			callbackErr = err
		}
	}

	if len(foundHWNDs) == 0 {
		if callbackErr != nil {
			return nil, fmt.Errorf("failed to find HWND for %s: %w", targetProcessName, callbackErr)
		}
		return nil, fmt.Errorf("could not find an open window for %s", targetProcessName)
	}

	return foundHWNDs, nil
}

func getDisplayDimensions() ([]displayDimensions, error) {
	user32DLL := syscall.NewLazyDLL("user32.dll")
	enumDisplayMonitorsProc := user32DLL.NewProc("EnumDisplayMonitors")
	getMonitorInfoWProc := user32DLL.NewProc("GetMonitorInfoW")

	displays := []displayDimensions{}
	displayNumber := 1

	callback := syscall.NewCallback(func(hMonitor uintptr, hdcMonitor uintptr, lprcMonitor uintptr, dwData uintptr) uintptr {
		info := monitorInfo{CbSize: uint32(unsafe.Sizeof(monitorInfo{}))}
		ret, _, _ := getMonitorInfoWProc.Call(
			hMonitor,
			uintptr(unsafe.Pointer(&info)),
		)
		if ret == 0 {
			return 1 // continue
		}

		width := int(info.RcMonitor.Right - info.RcMonitor.Left)
		height := int(info.RcMonitor.Bottom - info.RcMonitor.Top)

		displays = append(displays, displayDimensions{
			Height:        height,
			Width:         width,
			DisplayNumber: displayNumber,
			IsMain:        info.DwFlags&0x1 != 0, // MONITORINFOF_PRIMARY
			Left:          int(info.RcMonitor.Left),
			Top:           int(info.RcMonitor.Top),
			Right:         int(info.RcMonitor.Right),
			Bottom:        int(info.RcMonitor.Bottom),
		})
		displayNumber++

		return 1 // continue
	})

	ret, _, err := enumDisplayMonitorsProc.Call(
		0,
		0,
		callback,
		0,
	)
	if ret == 0 {
		return nil, fmt.Errorf("failed to enumerate displays: %w", err)
	}
	if len(displays) == 0 {
		return nil, fmt.Errorf("no displays found")
	}

	return displays, nil
}

func getWindowDimensions(validatedWindowName string) (map[string]int32, error) {
	user32DLL := syscall.NewLazyDLL("user32.dll")
	getWindowRectProc := user32DLL.NewProc("GetWindowRect")

	hwnds, err := getHWNDs(validatedWindowName)
	if err != nil {
		return nil, err
	}

	var data rect
	ret, _, err := getWindowRectProc.Call(hwnds[0], uintptr(unsafe.Pointer(&data)))
	if ret == 0 {
		return nil, err
	}

	dimensions := map[string]int32{}

	dimensions["left"] = data.Left
	dimensions["right"] = data.Right
	dimensions["top"] = data.Top
	dimensions["bottom"] = data.Bottom

	dimensions["x"] = data.Left
	dimensions["y"] = data.Top
	dimensions["width"] = data.Right - data.Left
	dimensions["height"] = data.Bottom - data.Top

	return dimensions, nil
}
