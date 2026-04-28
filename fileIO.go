package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileSpec struct {
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
	_, err = os.Stat(filepath.Join(wd, "place.json"))
	if err != nil {
		if os.IsNotExist(err) {
			// json doesn't exist, create it
			return fmt.Errorf("Not yet implemented: json doesn't exist, create it")
		}
		return err
	}

	// json does exist, or we just created it
	fmt.Println(name)
	return fmt.Errorf("Not yet implemented: createNewLocationAndSave")
}
