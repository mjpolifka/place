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
		case "locate":
			locate(args)
		case "track":
			if err := runTracker(); err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Println("Invalid argument, show help")
		}
	} else {
		fmt.Println("Not enough arguments, show help")
	}

}

func locate(args []string) {
	fmt.Println("Locate a window.")

	if len(args) < 8 {
		fmt.Println("Not enough args for 'locate', show help")
		return
	}

	// Process Name
	normalizedProcessName, err := normalizeProcessName(args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Process name:", normalizedProcessName)

	// Instance
	instance, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = validateIntOverflow(instance)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window instance:", instance)

	// X
	x, err := strconv.Atoi(args[4])
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
	y, err := strconv.Atoi(args[5])
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
	width, err := strconv.Atoi(args[6])
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
	height, err := strconv.Atoi(args[7])
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

	// Call Resize
	fmt.Println("Call resize")
	err = moveWindow(
		normalizedProcessName,
		instance,
		x,
		y,
		width,
		height,
	)
	if err != nil {
		fmt.Println(err)
	}
}
