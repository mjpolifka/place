//go:build windows

package main

import (
	"fmt"
	"syscall"
)

func moveWindow(windowName string, instance int, x int, y int, width int, height int) error {
	user32Dll := syscall.NewLazyDLL("user32.dll")
	moveWindowProc := user32Dll.NewProc("MoveWindow")
	setForegroundWindowProc := user32Dll.NewProc("SetForegroundWindow")
	bringWindowToTopProc := user32Dll.NewProc("BringWindowToTop")

	hwnd, err := getHWND(windowName)
	if err != nil {
		return err
	}

	fmt.Println("Not using instance:", instance)
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

	bringWindowToTopProc.Call(hwnd)
	setForegroundWindowProc.Call(hwnd)

	fmt.Println("Window moved, resized, and brought to the foreground.")
	return nil
}
