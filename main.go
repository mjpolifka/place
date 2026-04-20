package main

import (
	"fmt"
	"os"
	"strconv"
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
	err = moveWindow(
		args[2]+".exe",
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
