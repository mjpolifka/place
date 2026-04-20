//go:build !windows

package main

import "fmt"

func moveWindow(windowName string) error {
	return fmt.Errorf("moveWindow(%q) is only supported on Windows", windowName)
}
