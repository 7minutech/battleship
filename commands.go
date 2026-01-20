package main

import (
	"fmt"
	"os"
)

func runCommand(cmd string) {
	switch cmd {
	case "help":
		help()
	case "quit":
		quit()
	default:
		fmt.Printf("did not recognize command: %s\n", cmd)
	}
}

func help() {
	fmt.Println("help: prints possible commands")
	fmt.Println("quit: exits the program")
}

func quit() {
	fmt.Println("exiting program...")
	os.Exit(0)
}
