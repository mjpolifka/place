package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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

func createNewLocationAndSave(name string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath := filepath.Join(wd, "place.json")

	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// json doesn't exist, create it
			placeFile := PlaceFile{SelectedLocation: name, Locations: []Location{{Name: name, Places: []Place{}}}}
			if err := savePlaceFile(placeFile, filePath); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	// json does exist, or we just created it
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var placeFile PlaceFile
	if err = json.Unmarshal(fileBytes, &placeFile); err != nil {
		fmt.Println("Place.json is corrupt. Overwrite? y/N")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		choice := strings.TrimSpace(input)
		if choice == "y" || choice == "Y" {
			fmt.Println("Overwriting...")
			placeFile := PlaceFile{SelectedLocation: name, Locations: []Location{{Name: name, Places: []Place{}}}}
			if err := savePlaceFile(placeFile, filePath); err != nil {
				return err
			}
			return nil
		}
		fmt.Println("Exiting")
		return nil
	}
	fmt.Println(placeFile)
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
	newLocation := Location{Name: name, Places: []Place{}}
	placeFile.Locations = append(placeFile.Locations, newLocation)
	placeFile.SelectedLocation = name
	fmt.Println("New placeFile:", placeFile)
	savePlaceFile(placeFile, filePath)
	return nil
}

func savePlaceFile(placeFile PlaceFile, filePath string) error {
	jsonBytes, err := json.MarshalIndent(placeFile, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonBytes, 0644)
}
