//go:build windows

package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args

	if len(args) > 1 {
		switch args[1] {
		default:
			move(args)
		}
	} else {
		fmt.Println("Not enough arguments, show help")
	}

}

func move(args []string) {
	if len(args) < 6 {
		fmt.Println("Not enough args for 'move', show help")
		return
	}

	// Process Name
	normalizedProcessName, err := normalizeProcessName(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// X
	x, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Y
	y, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Width
	width, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Height
	height, err := strconv.Atoi(args[5])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Validate Dimensions
	err = validateDimensions(x, y, height, width)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Call Resize
	err = moveWindow(
		normalizedProcessName,
		x,
		y,
		width,
		height,
	)
	if err != nil {
		fmt.Println(err)
	}
}
