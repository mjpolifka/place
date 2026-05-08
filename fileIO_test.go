package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidatePlaceFile(t *testing.T) {
	// test file doesn't exist
	t.Run("file-doesnt-exist", func(t *testing.T) {
		// change to T.TempDir, then file won't exist
		t.Chdir(t.TempDir())

		// start test
		exist, _, _, err := validatePlaceFile()
		if err != nil {
			t.Fatal("want: exist==false | got: ERROR")
		}
		if exist {
			t.Error("want: exist==false | got: exist==true")
		}
	})
	// test file exists but is invalid
	t.Run("file-is-corrupt", func(t *testing.T) {
		// setup test
		tempDir := t.TempDir()
		t.Chdir(tempDir)
		path := filepath.Join(tempDir, "place.json")
		placeFile := `{"selected_location": "desktop","locations": [{"name": "laptop","places": ["bad"]},{"name": "desktop","places": []}]}`

		if err := os.WriteFile(path, []byte(placeFile), 0644); err != nil {
			t.Fatal("Couldn't set up test:", err)
		}

		// start test
		exist, valid, _, err := validatePlaceFile()
		if err != nil {
			t.Fatal("want: exist==true && valid==false | got:", err)
		}
		if !exist {
			t.Error("want: exist==true && valid==false | got: exist==false")
		}
		if valid {
			t.Error("want: exist==true && valid==false | got: valid==true")
		}
	})
	// test file exists and is valid
	t.Run("file-is-valid", func(t *testing.T) {
		// setup test
		tempDir := t.TempDir()
		t.Chdir(tempDir)
		path := filepath.Join(tempDir, "place.json")
		placeFile := PlaceFile{SelectedLocation: "desktop", Locations: []Location{}}
		placeFile.Locations = append(placeFile.Locations, Location{Name: "desktop", Places: []Place{}})
		jsonBytes, err := json.Marshal(placeFile)
		if err != nil {
			t.Fatal("Couldn't set up test:", err)
		}
		if err = os.WriteFile(path, jsonBytes, 0644); err != nil {
			t.Fatal("Couldn't set up test:", err)
		}

		// start test
		exist, valid, _, err := validatePlaceFile()
		if err != nil {
			t.Fatal("want: exist==true && valid==true | got:", err)
		}
		if !exist {
			t.Error("want: exist==true && valid==true | got: exist==false")
		}
		if !valid {
			t.Error("want: exist==true && valid==true | got: valid==false")
		}
	})
}

func TestGetUserInput(t *testing.T) {
	// test input string matches return string
	in := strings.NewReader("test string\n")
	got, err := getUserInput(in)
	if err != nil {
		t.Fatal("want: no error | got:", err)
	}
	if got != "test string" {
		t.Error("want: test string | got:", got)
	}
}

func TestAppendNewLocation(t *testing.T) {
	// test location already exists
	t.Run("test-already-exists", func(t *testing.T) {
		// setup test
		placeFile := PlaceFile{SelectedLocation: "desktop", Locations: []Location{}}
		placeFile.Locations = append(placeFile.Locations, Location{Name: "desktop", Places: []Place{}})

		// start test
		if err := appendNewLocation("desktop", &placeFile); err != nil {
			if err.Error() != "Can't create 'desktop', location already exists" {
				t.Error("want: Can't create 'desktop', location already exists | got:", err)
			}
		} else {
			t.Error("want: Can't create 'desktop', location already exists | got: PASS")
		}
	})
	// test location doesn't exist
	t.Run("test-doesnt-exist", func(t *testing.T) {
		// setup test
		placeFile := PlaceFile{SelectedLocation: "desktop", Locations: []Location{}}
		placeFile.Locations = append(placeFile.Locations, Location{Name: "desktop", Places: []Place{}})

		// start test
		if err := appendNewLocation("laptop", &placeFile); err != nil {
			t.Error("want: PASS | got:", err)
		}
	})
}

func TestSavePlaceFile(t *testing.T) {
	// test known good data
	t.Run("test-good-data", func(t *testing.T) {
		t.Error("Not yet implemented")
	})
	// test known bad data
	t.Run("test-bad-data", func(t *testing.T) {
		t.Error("Not yet implemented")
	})
}
