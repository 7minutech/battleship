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

	if err := pubsub.PublishJSON(ch, routing.EXCHANGE_BATTLESHIP_TOPIC, routing.NEW_PLAYER_KEY+"."+userName, gamelogic.NewPlayerMessage{UserName: userName}); err != nil {
		log.Fatalf("Failed to publish new player message: %v", err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.EXCHANGE_BATTLESHIP_TOPIC,
		routing.BOARD_STATE_KEY+"."+userName,
		routing.BOARD_STATE_KEY+".*",
		pubsub.Transient,
		gamelogic.ClientBoardStateHandler(userName),
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to board state messages: %v", err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.EXCHANGE_BATTLESHIP_TOPIC,
		routing.PLACE_STATE_KEY+"."+userName,
		routing.PLACE_STATE_KEY+".*",
		pubsub.Transient,
		gamelogic.ClientPlaceShipHandler(userName),
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to place ship messages: %v", err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.EXCHANGE_BATTLESHIP_TOPIC,
		routing.AUTO_PLACE_STATE_KEY+"."+userName,
		routing.AUTO_PLACE_STATE_KEY+".*",
		pubsub.Transient,
		gamelogic.ClientAutoPlaceHandler(userName),
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to auto place messages: %v", err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.EXCHANGE_BATTLESHIP_DIRECT,
		routing.GAME_RESET_KEY+"."+userName,
		routing.GAME_RESET_KEY,
		pubsub.Transient,
		gamelogic.ClientGameResetHandler,
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to game reset messages: %v", err)
	}

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">>> ")
		scanner.Scan()
		words := gamelogic.GetWords(scanner.Text())
		if words[0] == "help" {
			gamelogic.ClientHelp()
			continue
		} else if words[0] == "quit" {
			gamelogic.Quit()
			break
		} else if words[0] == "show" {
			pubsub.PublishJSON(ch, routing.EXCHANGE_BATTLESHIP_TOPIC, routing.SHOW_BOARD_KEY+"."+userName, routing.ShowBoardMessage{UserName: userName})
		} else if words[0] == "place" {
			if len(words) != 4 {
				fmt.Println("Usage: place <ship_type> <start_coord> <end_coord>")
				continue
			} else {
				shipType := words[1]
				startCoord := words[2]
				endCoord := words[3]
				placeShip := routing.PlaceShipCommand{UserName: userName, ShipType: shipType, StartCoord: startCoord, EndCoord: endCoord}
				pubsub.PublishJSON(ch, routing.EXCHANGE_BATTLESHIP_TOPIC, routing.PLACE_BOARD_KEY+"."+userName, placeShip)
			}
		} else if words[0] == "auto" {
			pubsub.PublishJSON(ch, routing.EXCHANGE_BATTLESHIP_TOPIC, routing.AUTO_PLACE_KEY+"."+userName, routing.AutoPlaceMessage{UserName: userName})
		} else if words[0] == "peek" {
			pubsub.PublishJSON(ch, routing.EXCHANGE_BATTLESHIP_TOPIC, routing.OPPONENT_BOARD_STATE_KEY+"."+userName, routing.ClientMessage{UserName: userName})
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
