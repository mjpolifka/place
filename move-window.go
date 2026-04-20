//go:build windows

package main

import (
	"fmt"
	"syscall"
)

func moveWindow() {
	user32Dll := syscall.NewLazyDLL("user32.dll")
	moveWindowProc := user32Dll.NewProc("MoveWindow")

	hwnd, err := getHWND()
	if err != nil {
		fmt.Println("Unable to find Notepad window:", err)
		return
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
		fmt.Println("MoveWindow failed:", err)
		return
	}

	fmt.Println("Window resized.")
}
