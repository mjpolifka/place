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
			if err := saveNewPlaceFile(name, filePath); err != nil {
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
			if err := saveNewPlaceFile(name, filePath); err != nil {
				return err
			}
			return nil
		}
		fmt.Println("Exiting")
		return nil
	}
	fmt.Println(placeFile)
	// check if name exists as a location
	return fmt.Errorf("Not yet implemented: createNewLocationAndSave")
}

func saveNewPlaceFile(name, filePath string) error {
	placeFile := PlaceFile{SelectedLocation: name, Locations: []Location{{Name: name, Places: []Place{}}}}
	jsonBytes, err := json.MarshalIndent(placeFile, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonBytes, 0644)
}
