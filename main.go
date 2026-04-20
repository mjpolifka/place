package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	args := os.Args

	// fmt.Println("All args:", args)

	if len(args) > 1 {
		if args[1] == "locate" {
			locate(args)
		} else {
			fmt.Println("Invalid argument, show help")
		}
	} else {
		fmt.Println("Not enough arguments, show help")
	}

}

func locate(args []string) {
	fmt.Println("Locate a window.")

	if len(args) < 8 {
		fmt.Println("Not enough args for 'locate', show help")
		return
	}

	// Process Name
	normalizedProcessName, err := normalizeProcessName(args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Process name:", normalizedProcessName)

	// Instance
	instance, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = validateIntOverflow(instance)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window instance:", instance)

	// X
	x, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = validateWindowCoord(x)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window x:", x)

	// Y
	y, err := strconv.Atoi(args[5])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = validateWindowCoord(y)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window y:", y)

	// Width
	width, err := strconv.Atoi(args[6])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = validateWindowSize(width)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window width:", width)

	// Height
	height, err := strconv.Atoi(args[7])
	if err != nil {
		fmt.Println(err)
		return
	}
	err = validateWindowSize(height)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window height:", height)

	// Call Resize
	fmt.Println("Call resize")
	err = moveWindow(
		normalizedProcessName,
		instance,
		x,
		y,
		width,
		height,
	)
	if err != nil {
		fmt.Println(err)
	}
}

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
