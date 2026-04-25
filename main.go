//go:build windows

package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args
	parseArgsAndRun(args)
}

func parseArgsAndRun(args []string) error {
	if len(args) > 1 {
		switch args[1] {
		default:
			err := move(args)
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("Not enough arguments, show help")
	}
	return nil
}

func move(args []string) error {
	if len(args) < 6 {
		return fmt.Errorf("Not enough args for 'move', show help")
	}

	// Process Name
	normalizedProcessName, err := normalizeProcessName(args[1])
	if err != nil {
		return err
	}

	// X
	x, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}

	// Y
	y, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}

	// Width
	width, err := strconv.Atoi(args[4])
	if err != nil {
		return err
	}

	// Height
	height, err := strconv.Atoi(args[5])
	if err != nil {
		return err
	}

	// Validate Dimensions
	err = validateDimensions(x, y, width, height)
	if err != nil {
		return err
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
		return err
	}

	return nil
}
