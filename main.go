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

	// TODO: validate input
	fmt.Println("Window name:", args[2])

	// TODO: validate input
	instance, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window instance:", instance)

	// TODO: validate input
	x, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window x:", x)

	// TODO: validate input
	y, err := strconv.Atoi(args[5])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window y:", y)

	// TODO: validate input
	width, err := strconv.Atoi(args[6])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window width:", width)

	// TODO: validate input
	height, err := strconv.Atoi(args[7])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Window height:", height)

	fmt.Println("Call resize")
	normalizedProcessName, err := normalizeProcessName(args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
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
