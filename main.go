//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
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

	wd := filepath.Dir(args[0])

	// Validate the input from the user before using it
	locationName := args[2]
	if err := validateLocationName(locationName); err != nil {
		return err
	}

	// Validate placeFile exists and is valid
	exist, valid, placeFile, err := validatePlaceFile(wd)
	if err != nil {
		return err
	}

	if !exist {
		// fmt.Println("Creating empty placeFile A")
		placeFile = PlaceFile{SelectedLocation: "", Locations: []Location{}}
	} else if !valid {
		fmt.Println("Place.json is corrupt. Overwrite? y/N")
		userInput, err := getUserInput(os.Stdin)
		if err != nil {
			return err
		}
		if userInput[0] == 'y' || userInput[0] == 'Y' {
			// fmt.Println("Creating empty placeFile B")
			placeFile = PlaceFile{SelectedLocation: "", Locations: []Location{}}
		} else {
			return nil
		} // exit
	}
	// only 3 ways out of the above block:
	// 	file exists and is valid,
	// 	file didn't exist and was created,
	// 	or file wasn't valid and was overwritten
	// no matter which path is taken, the file now exists and is valid, or we exited

	// exists and is valid, append new location if it doesn't exist
	if err = appendNewLocation(locationName, &placeFile); err != nil {
		return err
	}

	// fmt.Println("Saving placeFile")
	if err = savePlaceFile(wd, placeFile); err != nil {
		return err
	}

	return nil
}
