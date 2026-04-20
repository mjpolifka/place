package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args
	numArgs := len(args)
	// fmt.Println("All args:", args)

	var instance int
	var x int
	var y int
	var width int
	var height int
	var err error

	if len(args) > 1 {
		if args[1] == "locate" {
			fmt.Println("Locate a window.")

			if numArgs > 2 {
				fmt.Println("Window name:", args[2])
				// TODO: validate input

				if numArgs > 3 {
					fmt.Println("Window instance:", args[3])
					// TODO: validate input
					instance, err = strconv.Atoi(args[3])
					if err != nil {
						fmt.Println(err)
						os.Exit(0)
					}

					if numArgs > 4 {
						fmt.Println("Window x:", args[4])
						// TODO: validate input
						x, err = strconv.Atoi(args[4])
						if err != nil {
							fmt.Println(err)
							os.Exit(0)
						}

						if numArgs > 5 {
							fmt.Println("Window y:", args[5])
							// TODO: validate input
							y, err = strconv.Atoi(args[5])
							if err != nil {
								fmt.Println(err)
								os.Exit(0)
							}

							if numArgs > 6 {
								fmt.Println("Window width:", args[6])
								// TODO: validate input
								width, err = strconv.Atoi(args[6])
								if err != nil {
									fmt.Println(err)
									os.Exit(0)
								}

								if numArgs > 7 {
									fmt.Println("Window height:", args[7])
									// TODO: validate input
									height, err = strconv.Atoi(args[7])
									if err != nil {
										fmt.Println(err)
										os.Exit(0)
									}

									fmt.Println("Call resize")
									err := moveWindow(
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
									os.Exit(0)
								}
							}
						}
					}
				}
			}

			fmt.Println("Not enough args for 'locate', show help")
			os.Exit(0)
		}
	}

	fmt.Println("Invalid argument, show help")

	// err := moveWindow("notepad.exe")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
