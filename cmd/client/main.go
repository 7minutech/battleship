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
	player := gamelogic.CreatePlayer(userName)
	gameState := gamelogic.NewGameState(player)
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">>> ")
		scanner.Scan()
		words := gamelogic.GetWords(scanner.Text())
		switch words[0] {
		case "help":
			gamelogic.Help()
		case "place":
			if len(words) != 4 {
				log.Println("did not provide 3 args with place; usage: place cruiser a1 a5")
				continue
			}
			err := gameState.PlaceShip(words)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Printf("placed ship: %s from %s to %s", words[1], words[2], words[3])
		case "show":
			gameState.Show()
		case "quit":
			gamelogic.Quit()
		default:
			fmt.Printf("did not recognize command: %s\n", words[0])
		}
	}
}
