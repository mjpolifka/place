package main

import (
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

// func TestValidateIntOverflow(t *testing.T) {
// 	t.Error("Haven't implemented test")
// }

// func TestValidateDimensions(t *testing.T) {
// 	t.Error("Haven't implemented test")
// }
