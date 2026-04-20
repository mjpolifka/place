package main

import "fmt"

func main() {
	err := moveWindow("notepad.exe")
	if err != nil {
		fmt.Println(err)
	}
}
