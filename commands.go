package main

import (
	"fmt"
	"log"
	"os"
)

func runCommand(cmds []string) {
	switch cmds[0] {
	case "help":
		help()
	case "quit":
		quit()
	case "move":
		if len(cmds) != 2 {
			fmt.Println("move: needs one additonal word for location (move: a1)")
			return
		}
		log.Printf("moving to %s\n", cmds[1])
	default:
		fmt.Printf("did not recognize command: %s\n", cmds[0])
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
