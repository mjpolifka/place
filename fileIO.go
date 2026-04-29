package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
			jsonBytes, err := json.MarshalIndent(placeFile, "", "  ")
			if err != nil {
				return err
			}
			return os.WriteFile(filePath, jsonBytes, 0644)
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
		// ask to overwrite
		return err
	}
	fmt.Println(placeFile)
	// check if name exists as a location
	return fmt.Errorf("Not yet implemented: createNewLocationAndSave")
}
