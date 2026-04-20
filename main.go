package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	numArgs := len(args)
	// fmt.Println("All args:", args)

	if args[1] == "locate" {
		fmt.Println("Locate a window.")
		if numArgs > 2 {
			fmt.Println("Window name:", args[2])
			if numArgs > 3 {
				fmt.Println("Window instance:", args[3])
				if numArgs > 4 {
					fmt.Println("Window x:", args[4])
					if numArgs > 5 {
						fmt.Println("Window y:", args[5])
						if numArgs > 6 {
							fmt.Println("Window width:", args[6])
							if numArgs > 7 {
								fmt.Println("Window height:", args[7])
								fmt.Println("Call resize")
								os.Exit(0)
							}
						}
					}
				}
			}
		}

		fmt.Println("Not enough args")
		os.Exit(0)
	}

	fmt.Println("Invalid argument, show help")

	// err := moveWindow("notepad.exe")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
