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
	fmt.Println(name)
	return fmt.Errorf("Not yet implemented: createNewLocationAndSave")
}
