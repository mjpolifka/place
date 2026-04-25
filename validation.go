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
	maxWindowInt32 := 2147483647
	if i < -maxWindowInt32 || i > maxWindowInt32 {
		return fmt.Errorf("Int must be in signed 32-bit range [%d, %d]: %d", -maxWindowInt32, maxWindowInt32, i)
	}
	return nil
}

func validateDimensions(x, y, width, height int) error {
	// Check Int Overflow
	if err := validateIntOverflow(x); err != nil {
		return err
	}
	if err := validateIntOverflow(y); err != nil {
		return err
	}
	if err := validateIntOverflow(height); err != nil {
		return err
	}
	if err := validateIntOverflow(width); err != nil {
		return err
	}

	// Check against displays, first get display data
	displays, err := getDisplayDimensions()
	if err != nil {
		return err
	}

	// Check X, Y within displays
	foundDisplay := 0
	for _, display := range displays {
		if x >= display.Left && x <= display.Right {
			if y >= display.Top && y <= display.Bottom {
				foundDisplay = display.DisplayNumber
			}
		}
	}
	if foundDisplay == 0 {
		return fmt.Errorf("X, Y not within displays | x,y: %d,%d", x, y)
	}

	// Check height, width within display
	for _, display := range displays {
		if display.DisplayNumber == foundDisplay {
			if height <= display.Height && height >= 0 {
				if width <= display.Width && width >= 0 {
					return nil
				}
			}
		}
	}
	return fmt.Errorf("width, height larger than display | width,height: %d,%d", width, height)
}
