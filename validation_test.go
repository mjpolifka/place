package main

import (
	"strconv"
	"testing"
)

func TestNormalizeProcessName(t *testing.T) {
	t.Run("notepad", func(t *testing.T) {
		got, err := normalizeProcessName("NotepAd")
		if err != nil {
			t.Error(err)
		}
		if got != "notepad.exe" {
			t.Error("want: notepad.exe, got:", got)
		}
	})

	// test empty, expect err
	t.Run("empty", func(t *testing.T) {
		_, err := normalizeProcessName("")
		if err != nil {
			if err.Error() != "window name cannot be empty" {
				t.Error("want: window name cannot be empty, got:", err)
			}
			return
		}
		t.Error("want: FAIL, got: PASS")
	})

	// test only whitespace, expect err
	t.Run("whitespace", func(t *testing.T) {
		_, err := normalizeProcessName("    ")
		if err != nil {
			if err.Error() != "window name cannot be empty" {
				t.Error("want: window name cannot be empty, got:", err)
			}
			return
		}
		t.Error("want: FAIL, got: PASS")
	})

	// test with path separators, expect err
	t.Run("path-separator-forward", func(t *testing.T) {
		_, err := normalizeProcessName("fire/fox")
		if err != nil {
			if err.Error() != "window name cannot contain path separators" {
				t.Error("want: window name cannot contain path separators, got:", err)
			}
			return
		}
		t.Error("want: FAIL, got:PASS")
	})
	t.Run("path-separator-back", func(t *testing.T) {
		_, err := normalizeProcessName("fire\\fox")
		if err != nil {
			if err.Error() != "window name cannot contain path separators" {
				t.Error("want: window name cannot contain path separators, got:", err)
			}
			return
		}
		t.Error("want: FAIL, got: PASS")
	})

	// test with control characters, expect err
	t.Run("control-character", func(t *testing.T) {
		_, err := normalizeProcessName("fire\x00fox")
		if err != nil {
			if err.Error() != "window name cannot contain control characters" {
				t.Error("want: window name cannot contain control characters, got:", err)
			}
			return
		}
		t.Error("want: FAIL, got: PASS")
	})
}

func TestValidateIntOverflow(t *testing.T) {
	t.Run("good-int", func(t *testing.T) {
		err := validateIntOverflow(1024)
		if err != nil {
			t.Error("want: nil | got:", err)
		}
	})
	t.Run("bad-int-high", func(t *testing.T) {
		err := validateIntOverflow(2147483648)
		if err != nil {
			if err.Error() != "Int must be in signed 32-bit range [-2147483647, 2147483647]: 2147483648" {
				t.Error("want: Int must be in signed 32-bit range [-2147483647, 2147483647]: 2147483648 | got:", err)
			}
			return
		}
		t.Error("want: FAIL | got: PASS")
	})
	t.Run("bad-int-low", func(t *testing.T) {
		err := validateIntOverflow(-2147483648)
		if err != nil {
			if err.Error() != "Int must be in signed 32-bit range [-2147483647, 2147483647]: -2147483648" {
				t.Error("want: Int must be in signed 32-bit range [-2147483647, 2147483647]: -2147483648 | got:", err)
			}
			return
		}
		t.Error("want: FAIL | got: PASS")
	})
}

func TestValidateDimensions(t *testing.T) {
	displays, err := getDisplayDimensions()
	if err != nil {
		t.Error(err)
	}
	// test passing case
	t.Run("good-dimensions", func(t *testing.T) {
		width := 0
		height := 0
		for _, display := range displays {
			if display.Top == 0 && display.Left == 0 {
				width = display.Width
				height = display.Height
			}
		}
		if width == 0 || height == 0 {
			t.Error("Can't find display for 0,0")
		}
		err := validateDimensions(0, 0, width-1, height-1)
		if err != nil {
			t.Error("want: PASS | got:", err)
		}
	})
	// validateIntOverflow covered above
	// getDisplayDimensions covered in windows_calls_test
	// test x not within displays
	t.Run("x-not-within-displays", func(t *testing.T) {
		// find x outside displays
		right := 0
		for _, display := range displays {
			if display.Right > right {
				right = display.Right
			}
		}
		right += 1
		// then call validateDimensions with it
		err = validateDimensions(right, 0, 800, 800)
		if err != nil {
			if err.Error() != "X, Y not within displays | x,y: "+strconv.Itoa(right)+",0" {
				t.Error("want: X, Y not within displays | x,y: "+strconv.Itoa(right)+",0 | got:", err)
			}
			return
		}
		t.Error("want: FAIL | got: PASS")
	})
	// test y not within displays
	t.Run("y-not-within-displays", func(t *testing.T) {
		// find x outside displays
		bottom := 0
		for _, display := range displays {
			if display.Bottom > bottom {
				bottom = display.Bottom
			}
		}
		bottom += 1
		// then call validateDimensions with it
		err = validateDimensions(0, bottom, 800, 800)
		if err != nil {
			if err.Error() != "X, Y not within displays | x,y: 0,"+strconv.Itoa(bottom) {
				t.Error("want: X, Y not within displays | x,y: 0,"+strconv.Itoa(bottom)+" | got:", err)
			}
			return
		}
		t.Error("want: FAIL | got: PASS")
	})
	// test width larger than display
	t.Run("width-larger-than-displays", func(t *testing.T) {
		width := 0
		height := 0
		for _, display := range displays {
			if display.Top == 0 && display.Left == 0 {
				width = display.Width
				height = display.Height
			}
		}
		if width == 0 || height == 0 {
			t.Error("Can't find display for 0,0")
		}
		err = validateDimensions(0, 0, width+1, height-1)
		if err != nil {
			if err.Error() != "width, height larger than display | width,height: "+strconv.Itoa(width+1)+","+strconv.Itoa(height-1) {
				t.Error("want: width, height larger than display | width,height: "+strconv.Itoa(width+1)+","+strconv.Itoa(height-1)+" | got:", err)
			}
			return
		}
		t.Error("want: FAIL | got: PASS")
	})
	// test height larger than display
	t.Run("height-larger-than-displays", func(t *testing.T) {
		width := 0
		height := 0
		for _, display := range displays {
			if display.Top == 0 && display.Left == 0 {
				width = display.Width
				height = display.Height
			}
		}
		if width == 0 || height == 0 {
			t.Error("Can't find display for 0,0")
		}
		err = validateDimensions(0, 0, width-1, height+1)
		if err != nil {
			if err.Error() != "width, height larger than display | width,height: "+strconv.Itoa(width-1)+","+strconv.Itoa(height+1) {
				t.Error("want: width, height larger than display | width,height: "+strconv.Itoa(width-1)+","+strconv.Itoa(height+1)+" | got:", err)
			}
			return
		}
		t.Error("want: FAIL | got: PASS")
	})
}
