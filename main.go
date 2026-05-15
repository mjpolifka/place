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
		case "all":
			if len(args) > 2 {
				return fmt.Errorf("too many arguments for 'all'.  example: place all")
			}
			if err := all(wd); err != nil {
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
				if err := move(args); err != nil {
					return err
				}
				return nil
			}
			if len(args) == 2 {
				if err := defaultMove(wd, args[1]); err != nil {
					return err
				}
				return nil
			}
			return fmt.Errorf("could not parse arguments.  see readme.")
		}
	} else {
		return fmt.Errorf("Not enough arguments, show help")
	}
}

func move(args []string) error {
	if len(args) < 6 {
		return fmt.Errorf("Not enough args for 'move'.  Example: place firefox 0 0 1920 1080")
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

func defaultMove(wd string, processName string) error {
	normalizedProcessName, err := normalizeProcessName(processName)
	if err != nil {
		return err
	}
	placeFile, err := readPlaceFile(wd)
	if err != nil {
		return err
	}
	locationIndex, exists := placeFile.LocationMap()[placeFile.SelectedLocation]
	if !exists {
		return fmt.Errorf("currently selected location does not exist in place file: %s.  quitting.", placeFile.SelectedLocation)
	}
	placeIndex, exists := placeFile.Locations[locationIndex].PlaceMap()[normalizedProcessName]
	if !exists {
		return fmt.Errorf("no saved place for %s.  quitting.", normalizedProcessName)
	}
	selectedPlace := placeFile.Locations[locationIndex].Places[placeIndex]
	if err := moveWindow(
		selectedPlace.Name,
		selectedPlace.X,
		selectedPlace.Y,
		selectedPlace.Width,
		selectedPlace.Height,
	); err != nil {
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
	locationIndex, exists := placeFile.LocationMap()[placeFile.SelectedLocation]
	if !exists {
		return fmt.Errorf("selected location does not exist! quitting.")
	}
	placeIndex, exists := placeFile.Locations[locationIndex].PlaceMap()[normalizedProcessName]
	if !exists {
		// create a new Place object and fill in the deets
		place := Place{
			Name:   normalizedProcessName,
			X:      int(dimensions["x"]),
			Y:      int(dimensions["y"]),
			Width:  int(dimensions["width"]),
			Height: int(dimensions["height"]),
		}
		// append it to the right place
		placeFile.Locations[locationIndex].Places = append(placeFile.Locations[locationIndex].Places, place)
	} else {
		// edit the existing matching Place object with the deets
		placeFile.Locations[locationIndex].Places[placeIndex] = Place{
			Name:   normalizedProcessName,
			X:      int(dimensions["x"]),
			Y:      int(dimensions["y"]),
			Width:  int(dimensions["width"]),
			Height: int(dimensions["height"]),
		}
	}

	// then save
	savePlaceFile(wd, placeFile)
	fmt.Printf("saved new place for %s\n", normalizedProcessName)
	return nil
}

func all(wd string) error {
	placeFile, err := readPlaceFile(wd)
	if err != nil {
		return err
	}
	locationIndex, exists := placeFile.LocationMap()[placeFile.SelectedLocation]
	if !exists {
		return fmt.Errorf("selected location does not exist in place file: %s", placeFile.SelectedLocation)
	}

	for _, place := range placeFile.Locations[locationIndex].Places {
		if err := moveWindow(
			place.Name,
			place.X,
			place.Y,
			place.Width,
			place.Height,
		); err != nil {
			return err
		}
	}

	return nil
}
