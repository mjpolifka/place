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
	wd := filepath.Dir(args[0])
	if len(args) > 1 {
		switch args[1] {
		case "create":
			if len(args) > 3 {
				return fmt.Errorf("Too many args for 'create'")
			}
			if len(args) < 3 {
				return fmt.Errorf("Not enough args for 'create'")
			}
			locationName := args[2]
			if err := validateLocationName(locationName); err != nil {
				return err
			} else if err := create(wd, locationName); err != nil {
				return err
			}
			return nil
		case "select":
			if len(args) > 3 {
				return fmt.Errorf("Too many args for 'select'")
			}
			if len(args) < 3 {
				return fmt.Errorf("Not enough args for 'select'")
			}
			locationName := args[2]
			if err := validateLocationName(locationName); err != nil {
				return err
			} else if err := selectLocation(wd, locationName); err != nil {
				return err
			}
			return nil
		case "save":
			if len(args) > 3 {
				return fmt.Errorf("Too many args for 'save'")
			}
			if len(args) < 3 {
				return fmt.Errorf("Not enough args for 'save'")
			}
			processName := args[2]
			if err := save(wd, processName); err != nil {
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

func create(wd string, locationName string) error {
	placeFile, err := validatePlaceFile(wd)
	if err != nil {
		return err
	}

	// place.json exists and is valid, append new location if the location doesn't exist
	if err = appendNewLocation(locationName, &placeFile); err != nil {
		return err
	}

	// fmt.Println("Saving placeFile")
	if err = savePlaceFile(wd, placeFile); err != nil {
		return err
	}

	fmt.Printf("%s created\n", locationName)

	return nil
}

// "select" is a reserved word, hence the strange name
func selectLocation(wd string, locationName string) error {
	placeFile, err := validatePlaceFile(wd)
	if err != nil {
		return err
	}

	// Check if the location exists already
	exists := false
	for _, location := range placeFile.Locations {
		if location.Name == locationName {
			exists = true
		}
	}

	// If it doesn't exist, ask to create it
	if !exists {
		fmt.Printf("Location '%s' doesn't exist, create it? y/N", locationName)
		choice, err := getUserInput(os.Stdin)
		if err != nil {
			return err
		}
		if choice[0] == 'y' || choice[0] == 'Y' {
			// If the user said yes, create it
			appendNewLocation(locationName, &placeFile)
			fmt.Printf("%s created\n", locationName)
		} else {
			os.Exit(0)
		}
	}

	// Set the SelectedLocation to the new location and save
	placeFile.SelectedLocation = locationName
	savePlaceFile(wd, placeFile)
	fmt.Printf("%s selected\n", locationName)
	return nil
}

func save(wd string, processName string) error {
	normalizedProcessName, err := normalizeProcessName(processName)
	if err != nil {
		return err
	}
	dimensions, err := getWindowDimensions(normalizedProcessName)
	if err != nil {
		return err
	}

	placeFile, err := readPlaceFile(wd)
	if err != nil {
		return err
	}

	fmt.Println(dimensions["width"])
	fmt.Println(placeFile.SelectedLocation)
	return nil
}
