package main

import (
	"strings"
	"testing"
)

func TestValidatePlaceFile(t *testing.T) {
	// test file doesn't exist
	// test file exists but is invalid
	// test file exists and is valid
	t.Error("Not yet implemented")
}

func TestGetUserInput(t *testing.T) {
	// test input string matches return string
	in := strings.NewReader("test string\n")
	got, err := getUserInput(in)
	if err != nil {
		t.Error("want: no error | got:", err)
	}
	if got != "test string" {
		t.Error("want: test string | got:", got)
	}
}

func TestAppendNewLocation(t *testing.T) {
	// test location already exists
	// test location doesn't exist
	t.Error("Not yet implemented")
}

func TestSavePlaceFile(t *testing.T) {
	// test known good data
	// test known bad data
	t.Error("Not yet implemented")
}
