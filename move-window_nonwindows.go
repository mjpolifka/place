//go:build !windows

package main

import "fmt"

func moveWindow() {
	fmt.Println("moveWindow is only supported on Windows.")
}
