package main

import "testing"

func TestMoveWindow(t *testing.T) {
	t.Skip("Not implementing yet, not sure how")
}

func TestGetHWNDs(t *testing.T) {
	t.Skip("Not implementing yet, not sure how")
}

func TestGetDisplayDimensions(t *testing.T) {
	// test no errors
	_, err := getDisplayDimensions()
	if err != nil {
		t.Error("want: PASS | got:", err)
	}
}
