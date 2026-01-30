package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/7minutech/battleship/internal/gamelogic"
	"github.com/7minutech/battleship/internal/pubsub"
	"github.com/7minutech/battleship/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	fmt.Println("Connected to rabbit mq")
	defer conn.Close()

	userName, err := gamelogic.Welcome()
	if err != nil {
		log.Fatalf("error: could not get username %v", err)
	}
	fmt.Println("welcome", userName)
	player := gamelogic.CreatePlayer(userName)
	gameState := gamelogic.NewGameState(player)
	gameState.Show()

	err = pubsub.SubscribeJSON(
		conn,
		routing.EXCHANGE_BATTLESHIP_DIRECT,
		routing.PAUSE_KEY+"."+userName, routing.PAUSE_KEY,
		pubsub.Durabale,
		gamelogic.PauseHandler(gameState),
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to pause messages: %v", err)
	}

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">>> ")
		scanner.Scan()
		words := gamelogic.GetWords(scanner.Text())
		if words[0] == "help" {
			gamelogic.Help()
			continue
		} else if words[0] == "quit" {
			gamelogic.Quit()
			break
		} else if words[0] == "place" {
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
			gameState.Show()
		} else if words[0] == "show" {
			gameState.Show()
		} else if words[0] == "look" {
			gameState.ShowOpponentBoard()
		} else {
			fmt.Printf("did not recognize command: %s\n", words[0])
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")
}
