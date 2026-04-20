package main

import (
	"fmt"
	"syscall"
)

func user32() {
	user32 := syscall.NewLazyDLL("user32.dll") // is this a real DLL?
	moveWindow := user32.NewProc("MoveWindow")

	// Replace this with the HWND you want to resize.
	var hwnd uintptr = 0x00123456 // what's an HWND and how do I get one?

	x := 100
	y := 100
	width := 800
	height := 600
	repaint := 1 // TRUE
	// Do we need to choose a display?

	ret, _, err := moveWindow.Call(
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
