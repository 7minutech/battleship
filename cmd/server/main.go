package main

import (
	"bufio"
	"fmt"
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
		fmt.Printf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Failed to open a channel: %v", err)
		return
	}

	gameState := gamelogic.NewGameState()

	pubsub.SubscribeJSON(
		conn,
		routing.EXCHANGE_BATTLESHIP_DIRECT,
		routing.NEW_PLAYER_KEY+"."+"notifier",
		routing.NEW_PLAYER_KEY,
		pubsub.Durabale,
		gamelogic.NewPlayerHandler(gameState))

	gamelogic.PrintServerHelp()

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
		} else if words[0] == "pause" {
			pubsub.PublishJSON(ch, routing.EXCHANGE_BATTLESHIP_DIRECT, routing.PAUSE_KEY, routing.PauseMessage{Content: "Server is paused"})
			continue
		} else {
			fmt.Printf("did not recognize command: %s\n", words[0])
			continue
		}
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")
}
