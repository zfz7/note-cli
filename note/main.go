package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 1 && len(os.Args) != 2 {
		fmt.Println("Invalid number of arguments, accept 0 or 1")
		os.Exit(1)
	}

	if len(os.Args) == 2 && os.Args[1] == "setup" {
		err := WriteDefaultConfig()
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	relativeWeek, err := ReadRelativeWeek()
	if err != nil {
		os.Exit(1)
	}
	config, err := ReadConfig()
	if err != nil {
		os.Exit(1)
	}
	err = OpenNote(relativeWeek, config)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
