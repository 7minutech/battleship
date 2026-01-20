package gamelogic

import (
	"fmt"
	"log"
	"os"
)

func gameCommand(cmds []string) {
	switch cmds[0] {
	case "help":
		help()
	case "quit":
		quit()
	case "place":
		if len(cmds) != 4 {
			fmt.Println("place: needs 4 args (place s1 a1 a3)")
			return
		}
		log.Printf("placeing ship: %s at start: %s end: %s\n", cmds[1], cmds[2], cmds[3])
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
