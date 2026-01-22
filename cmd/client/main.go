package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/7minutech/battleship/internal/gamelogic"
)

func main() {
	userName, err := gamelogic.Welcome()
	if err != nil {
		log.Fatalf("error: could not get username %v", err)
	}
	fmt.Println("welcome", userName)
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">>> ")
		scanner.Scan()
		words := gamelogic.GetWords(scanner.Text())
		fmt.Print(words)
		switch words[0] {
		case "help":
			gamelogic.Help()
		case "quit":
			gamelogic.Quit()
		default:
			fmt.Printf("did not recognize command: %s\n", words[0])
		}
	}
}
