package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type PlaceFile struct {
	SelectedLocation string     `json:"selected_location"`
	Locations        []Location `json:"locations"`
}

type Location struct {
	Name   string  `json:"name"`
	Places []Place `json:"places"`
}

type Place struct {
	Name   string `json:"name"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type InvalidPlaceFileError struct {
	Err error
}

func (e *InvalidPlaceFileError) Error() string {
	return "place file is invalid: " + e.Err.Error()
}

func (e *InvalidPlaceFileError) Unwrap() error {
	return e.Err
}

func (placeFile PlaceFile) LocationMap() map[string]int {
	locMap := map[string]int{}
	for i, location := range placeFile.Locations {
		locMap[location.Name] = i
	}
	return locMap
}

func IsInvalidPlaceFile(err error) bool {
	var invalidErr *InvalidPlaceFileError
	return errors.As(err, &invalidErr)
}

func validatePlaceFile(wd string) (PlaceFile, error) {
	placeFile, err := readPlaceFile(wd)
	if err != nil {
		if os.IsNotExist(err) {
			// fmt.Println("Creating empty placeFile A")
			placeFile = PlaceFile{SelectedLocation: "", Locations: []Location{}}
			return placeFile, nil
		} else if IsInvalidPlaceFile(err) {
			fmt.Println("Place.json is corrupt. Overwrite? y/N")
			userInput, err := getUserInput(os.Stdin)
			if err != nil {
				return PlaceFile{}, err
			}
			if userInput[0] == 'y' || userInput[0] == 'Y' {
				// fmt.Println("Creating empty placeFile B")
				placeFile = PlaceFile{SelectedLocation: "", Locations: []Location{}}
				return placeFile, nil
			} else {
				os.Exit(0)
			} // exit
		} else {
			return PlaceFile{}, err
		}
	}
	return placeFile, nil
}

func readPlaceFile(wd string) (PlaceFile, error) { // exist, valid, placeFile, err
	filePath := filepath.Join(wd, "place.json")

	// check if json exists
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return PlaceFile{}, err
	}
	// json does exist
	// fmt.Println("JSON does exist in readPlaceFile")
	// fmt.Println("fileBytes:", string(fileBytes))

	// check if json is valid
	var placeFile PlaceFile
	if err = json.Unmarshal(fileBytes, &placeFile); err != nil {
		// json is not valid
		// fmt.Println("JSON is not valid, returning from readPlaceFile")
		return PlaceFile{}, &InvalidPlaceFileError{Err: err}
	}
	//json is valid
	// fmt.Println("JSON exists and is valid, returning from readPlaceFile")
	return placeFile, nil
}

func getUserInput(in io.Reader) (string, error) {
	reader := bufio.NewReader(in)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return "N", nil
	} else {
		return input, nil
	}
}

func appendNewLocation(name string, placeFile *PlaceFile) error {
	// check if name exists as a location
	found := false
	for _, location := range placeFile.Locations {
		if location.Name == name {
			found = true
		}
	}
	if found {
		return fmt.Errorf("Can't create '%s', location already exists", name)
	}

	// name doesn't exist, append it to existing
	// fmt.Println("Name doesn't exist, appending to existing file")
	newLocation := Location{Name: name, Places: []Place{}}
	placeFile.Locations = append(placeFile.Locations, newLocation)
	placeFile.SelectedLocation = name

	return nil
}

func savePlaceFile(wd string, placeFile PlaceFile) error {
	filePath := filepath.Join(wd, "place.json")

	jsonBytes, err := json.MarshalIndent(placeFile, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonBytes, 0644)
}
