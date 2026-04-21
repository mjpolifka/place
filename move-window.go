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
	showWindowProc := user32Dll.NewProc("ShowWindow")

	hwndList, err := getHWNDs(windowName)
	if err != nil {
		return err
	}

	fmt.Println("HWND List:", hwndList)

	fmt.Println("Not using instance:", instance)
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
	windows, err := enumerateOpenWindows()
	if err != nil {
		return nil, err
	}

	var foundHWNDs []uintptr
	for _, w := range windows {
		if w.ProcessName == targetProcessName {
			foundHWNDs = append(foundHWNDs, w.HWND)
		}
	}

	if len(foundHWNDs) == 0 {
		return nil, fmt.Errorf("could not find an open window for %s", targetProcessName)
	}

	return sortHWNDsByTrackedOrder(foundHWNDs), nil
}
