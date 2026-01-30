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

	gamelogic.PrintServerHelp()

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">>> ")
		scanner.Scan()
		words := gamelogic.GetWords(scanner.Text())
		switch words[0] {
		case "help":
			gamelogic.Help()
		case "pause":
			pubsub.PublishJSON(ch, routing.EXCHANGE_BATTLESHIP_DIRECT, routing.PAUSE_KEY, routing.PauseMessage{Content: "Server is paused"})
		case "quit":
			gamelogic.Quit()
		default:
			fmt.Printf("did not recognize command: %s\n", words[0])
		}

		// wait for ctrl+c
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		<-signalChan
		fmt.Println("RabbitMQ connection closed.")
	}
}
