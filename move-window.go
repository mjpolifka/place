//go:build windows

package main

import (
	"fmt"
	"syscall"
)

func moveWindow(windowName string) error {
	user32Dll := syscall.NewLazyDLL("user32.dll")
	moveWindowProc := user32Dll.NewProc("MoveWindow")

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

	fmt.Println("Window resized.")
	return nil
}
