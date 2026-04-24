//go:build windows

package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

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

	fmt.Println("HWND List:", hwndList)
	repaint := 1 // TRUE

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

	fmt.Println("Window restored, moved, resized, and brought to the foreground.")
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
