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

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	userName, err := gamelogic.Welcome()
	if err != nil {
		log.Fatalf("error: could not get username %v", err)
	}
	fmt.Println("welcome", userName)

	if err := pubsub.PublishJSON(ch, "battleship_direct", routing.NEW_PLAYER_KEY, gamelogic.NewPlayerMessage{UserName: userName}); err != nil {
		log.Fatalf("Failed to publish new player message: %v", err)
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
		} else {
			fmt.Printf("did not recognize command: %s\n", words[0])
			continue
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")
}
