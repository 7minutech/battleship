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

	err = pubsub.SubscribeJSON(
		conn,
		routing.EXCHANGE_BATTLESHIP_TOPIC,
		routing.NEW_PLAYER_KEY+"."+"notifier",
		routing.NEW_PLAYER_KEY+".*",
		pubsub.Transient,
		gamelogic.NewPlayerHandler(gameState),
	)
	if err != nil {
		fmt.Printf("Failed to subscribe to new player messages: %v", err)
		return
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.EXCHANGE_BATTLESHIP_TOPIC,
		routing.GAME_COMMANDS_KEY+"."+"show_board",
		routing.SHOW_BOARD_KEY+".*",
		pubsub.Transient,
		gamelogic.ShowBoardHandler(gameState, ch),
	)
	if err != nil {
		fmt.Printf("Failed to subscribe to show board messages: %v", err)
		return
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.EXCHANGE_BATTLESHIP_TOPIC,
		routing.GAME_COMMANDS_KEY+"."+"place_ship",
		routing.PLACE_BOARD_KEY+".*",
		pubsub.Transient,
		gamelogic.PlaceShipHandler(gameState, ch),
	)
	if err != nil {
		fmt.Printf("Failed to subscribe to place ship messages: %v", err)
		return
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.EXCHANGE_BATTLESHIP_TOPIC,
		routing.GAME_COMMANDS_KEY+"."+"auto_place",
		routing.AUTO_PLACE_KEY+".*",
		pubsub.Transient,
		gamelogic.AutoPlaceHandler(gameState, ch),
	)
	if err != nil {
		fmt.Printf("Failed to subscribe to auto place messages: %v", err)
		return
	}

	gamelogic.PrintServerHelp()

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">>> ")
		scanner.Scan()
		words := gamelogic.GetWords(scanner.Text())

		if words[0] == "help" {
			gamelogic.ServerHelp()
			continue
		} else if words[0] == "quit" {
			gamelogic.Quit()
			break
		} else if words[0] == "pause" {
			err := pubsub.PublishJSON(ch, routing.EXCHANGE_BATTLESHIP_DIRECT, routing.PAUSE_KEY, routing.PauseMessage{Content: "Server is paused"})
			if err != nil {
				fmt.Printf("Failed to publish pause message: %v", err)
			}
			continue
		} else if words[0] == "reset" {
			gameState.ResetGame()
			err := pubsub.PublishJSON(ch, routing.EXCHANGE_BATTLESHIP_DIRECT, routing.GAME_RESET_KEY, routing.GameResetMessage{Content: "Game has been reset"})
			if err != nil {
				fmt.Printf("Failed to publish game reset message: %v", err)
			}
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
