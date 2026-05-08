package main

import (
	"bufio"
	"encoding/json"
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

func validatePlaceFile(wd string) (bool, bool, PlaceFile, error) { // exist, valid, placeFile, err
	filePath := filepath.Join(wd, "place.json")

	// check if json exists
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// json doesn't exist
			// fmt.Println("JSON doesn't exist, returning from validatePlaceFile")
			return false, false, PlaceFile{}, nil
		}
		return false, false, PlaceFile{}, err
	}
	// json does exist
	// fmt.Println("JSON does exist in validatePlaceFile")
	// fmt.Println("fileBytes:", string(fileBytes))

	// check if json is valid
	var placeFile PlaceFile
	if err = json.Unmarshal(fileBytes, &placeFile); err != nil {
		// json is not valid
		// fmt.Println("JSON is not valid, returning from validatePlaceFile")
		return true, false, PlaceFile{}, nil // error == nil because this is how we check valid == false
	}
	//json is valid
	// fmt.Println("JSON exists and is valid, returning from validatePlaceFile")
	return true, true, placeFile, nil
}

func getUserInput(in io.Reader) (string, error) {
	reader := bufio.NewReader(in)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
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
