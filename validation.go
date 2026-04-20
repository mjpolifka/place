//go:build windows

package main

import (
	"fmt"
	"strings"
	"unicode"
)

func normalizeProcessName(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("window name cannot be empty")
	}

	for _, ch := range trimmed {
		if ch == '/' || ch == '\\' {
			return "", fmt.Errorf("window name cannot contain path separators")
		}
		if unicode.IsControl(ch) {
			return "", fmt.Errorf("window name cannot contain control characters")
		}
	}

	normalized := strings.ToLower(trimmed)
	normalized = strings.TrimSuffix(normalized, ".exe")
	normalized = normalized + ".exe"

	return normalized, nil
}

func validateIntOverflow(i int) error {
	minWindowInt32 := -2147483648
	maxWindowInt32 := 2147483647
	if i < minWindowInt32 || i > maxWindowInt32 {
		return fmt.Errorf("Int must be in signed 32-bit range [%d, %d]: %d", minWindowInt32, maxWindowInt32, i)
	}
	return nil
}

func validateWindowCoord(value int) error {
	maxWindowCoord := 10000
	err := validateIntOverflow(value)
	if err != nil {
		return err
	}
	if value < 0 || value > maxWindowCoord {
		return fmt.Errorf("Window coord must be in range [0, %d]: %d", maxWindowCoord, value)
	}
	return nil
}

func validateWindowSize(value int) error {
	maxWindowSize := 10000
	err := validateIntOverflow(value)
	if err != nil {
		return err
	}
	if value < 0 || value > maxWindowSize {
		return fmt.Errorf("Window size must be in range [0, %d]: %d", maxWindowSize, value)
	}
	return nil
}
