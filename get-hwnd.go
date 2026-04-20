//go:build windows

package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

func getHWND(targetProcessName string) (uintptr, error) {
	user32DLL := syscall.NewLazyDLL("user32.dll")
	kernel32DLL := syscall.NewLazyDLL("kernel32.dll")

	enumWindowsProc := user32DLL.NewProc("EnumWindows")
	getWindowThreadProcessIDProc := user32DLL.NewProc("GetWindowThreadProcessId")
	isWindowVisibleProc := user32DLL.NewProc("IsWindowVisible")

	openProcessProc := kernel32DLL.NewProc("OpenProcess")
	queryFullProcessImageNameWProc := kernel32DLL.NewProc("QueryFullProcessImageNameW")
	closeHandleProc := kernel32DLL.NewProc("CloseHandle")

	const processQueryLimitedInformation = 0x1000
	var foundHWND uintptr
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
			foundHWND = hwnd
			return 0 // stop enumeration
		}

		return 1 // continue
	})

	ret, _, err := enumWindowsProc.Call(callback, 0)
	if ret == 0 && foundHWND == 0 {
		if err != syscall.Errno(0) {
			callbackErr = err
		}
	}

	if foundHWND == 0 {
		if callbackErr != nil {
			return 0, fmt.Errorf("failed to find HWND for %s: %w", targetProcessName, callbackErr)
		}
		return 0, fmt.Errorf("could not find an open window for %s", targetProcessName)
	}

	return foundHWND, nil
}
