//go:build windows

package main

import (
	"fmt"
	"syscall"
)

func moveWindow(windowName string) error {
	user32Dll := syscall.NewLazyDLL("user32.dll")
	moveWindowProc := user32Dll.NewProc("MoveWindow")
	showWindowProc := user32Dll.NewProc("ShowWindow")
	setForegroundWindowProc := user32Dll.NewProc("SetForegroundWindow")
	bringWindowToTopProc := user32Dll.NewProc("BringWindowToTop")

	hwnd, err := getHWND(windowName)
	if err != nil {
		return err
	}

	x := 100
	y := 100
	width := 800
	height := 600
	repaint := 1 // TRUE

	ret, _, err := moveWindowProc.Call(
		hwnd,
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(repaint),
	)
	if ret == 0 {
		return err
	}

	const swRestore = 9
	showWindowProc.Call(hwnd, uintptr(swRestore))
	bringWindowToTopProc.Call(hwnd)
	setForegroundWindowProc.Call(hwnd)

	fmt.Println("Window moved, resized, and brought to the foreground.")
	return nil
}
