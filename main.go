//go:build windows

package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args

	// fmt.Println("All args:", args)

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
	fmt.Println("Move a window.")

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
	fmt.Println("Process name:", normalizedProcessName)

	// X
	x, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = validateWindowCoord(x)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window x:", x)

	// Y
	y, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = validateWindowCoord(y)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window y:", y)

	// Width
	width, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = validateWindowSize(width)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window width:", width)

	// Height
	height, err := strconv.Atoi(args[5])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = validateWindowSize(height)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window height:", height)

	displays, err := getDisplayDimensions()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = validatePointWithinDisplays(x, y, displays)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Call Resize
	fmt.Println("Call resize")
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
