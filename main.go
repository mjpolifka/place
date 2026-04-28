//go:build windows

package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args
	if err := parseArgsAndRun(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseArgsAndRun(args []string) error {
	if len(args) > 1 {
		switch args[1] {
		case "create":
			err := create(args)
			if err != nil {
				return err
			}
			return nil
		default:
			if len(args) > 2 {
				if args[2] == "is" {
					err := is(args)
					if err != nil {
						return err
					}
					return nil
				}
			}
			err := move(args)
			if err != nil {
				return err
			}
			return nil
		}
	} else {
		return fmt.Errorf("Not enough arguments, show help")
	}
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

func is(args []string) error {
	normalizedProcessName, err := normalizeProcessName(args[1])
	if err != nil {
		return err
	}
	data, err := getWindowDimensions(normalizedProcessName)
	if err != nil {
		return err
	}
	fmt.Printf("%s: %d %d %d %d\n", normalizedProcessName, data["x"], data["y"], data["width"], data["height"])
	return nil
}

func create(args []string) error {
	if len(args) > 3 {
		return fmt.Errorf("Too many args for 'create'")
	}
	if len(args) < 3 {
		return fmt.Errorf("Not enough args for 'create'")
	}

	locationName := args[2]
	if err := validateLocationName(locationName); err != nil {
		return err
	}
	fmt.Println("Location is valid:", locationName)
	return fmt.Errorf("Haven't implemented yet")
}
