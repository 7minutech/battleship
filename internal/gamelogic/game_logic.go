package gamelogic

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/7minutech/battleship/internal/pubsub"
	"github.com/7minutech/battleship/internal/routing"
	"github.com/olekukonko/tablewriter"
	amqp "github.com/rabbitmq/amqp091-go"
)

func GetWords(input string) []string {
	words := strings.Split(input, " ")
	return words
}

func convertMove(mv string) (boardMove, error) {
	if len(mv) > 3 {
		return boardMove{}, fmt.Errorf("error: move is too long: %s", mv)
	}

	letter := mv[0]
	num := string(mv[1:])

	numInt, err := strconv.Atoi(num)
	if err != nil {
		return boardMove{}, fmt.Errorf("error: could not convert number in move to int: %s", mv)
	}

	if !(int(START_COL_HEADER) <= int(letter) && int(letter) <= int(END_COL_HEADER)) {
		return boardMove{}, fmt.Errorf("error: move did not contain a letter that was between a-j: %s", mv)
	}

	if !(START_ROW_HEADER <= numInt && numInt <= END_ROW_HEADER) {
		return boardMove{}, fmt.Errorf("error: move did not contain a number that was between 1-10: %s", mv)
	}

	row := (int(letter) - int(START_COL_HEADER))
	col := numInt - START_ROW_HEADER

	return boardMove{row: row, col: col}, nil
}

func CreatePlayer(name string) Player {
	player := Player{
		userName:  name,
		shipCount: STARTING_SHIP_COUNT,
		ships:     createShips(),
	}

	return player
}

func createShips() []ship {
	ships := []ship{
		startCarrier,
		startBattleship,
		startCruiser,
		startDestroyer,
		startSubmarine,
	}

	return ships
}

func (gs *gameState) PlaceShip(words []string, board board) error {
	sp, err := gs.getShipPlacement(words, *gs.getPlayerByName(board.owner))
	if err != nil {
		return err
	}
	if err := gs.validateShipPlacement(sp, board); err != nil {
		return err
	}
	for i := range sp.ship.length {
		row := sp.start.row
		col := sp.start.col
		if sp.orientation == horizantal {
			col = sp.start.col + i
		} else {
			row = sp.start.row + i
		}
		bm := boardMove{row: row, col: col}
		board.sqaures[row][col] = &sp.ship
		sp.ship.modules[bm] = false
	}

	return nil
}

func (gs *gameState) getShipPlacement(words []string, player Player) (shipPlacement, error) {
	shipName := words[1]
	startPlace := words[2]
	endPlace := words[3]
	ship, err := gs.getShip(shipName, player)
	if err != nil {
		return shipPlacement{}, err
	}

	if len(ship.modules) != 0 {
		return shipPlacement{}, fmt.Errorf("error: %s has already been placed", shipName)
	}

	startPlaceMove, err := convertMove(startPlace)
	if err != nil {
		return shipPlacement{}, err
	}

	endPlaceMove, err := convertMove(endPlace)
	if err != nil {
		return shipPlacement{}, err
	}

	var shipOrientation orientation
	if startPlaceMove.row == endPlaceMove.row {
		shipOrientation = horizantal
	} else if startPlaceMove.col == endPlaceMove.col {
		shipOrientation = vertical
	} else {
		return shipPlacement{}, fmt.Errorf("error: ship start and end were not on same row or col")
	}

	shipPlacement := shipPlacement{start: startPlaceMove, end: endPlaceMove, ship: ship, orientation: shipOrientation}

	return shipPlacement, err

}

func (gs *gameState) validateShipPlacement(sp shipPlacement, board board) error {
	var err error
	err = validShipRange(sp)
	err = gs.shipsOccupyRange(sp, board)
	return err
}

var ErrShipNotFound = errors.New("error: could not find ship with that name")

func (gs *gameState) getShip(shipName string, player Player) (ship, error) {
	for _, ship := range player.ships {
		if ship.name == shipName {
			return ship, nil
		}
	}
	return ship{}, ErrShipNotFound
}

func Show(boardData [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header(boardData[0])
	table.Bulk(boardData[1:])
	table.Render()

}

func (gs *gameState) boardData(board board) [][]string {
	var data [][]string
	header := []string{" ", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	data = append(data, header)
	for row := range BOARD_SIZE {
		rowLabelVal := A_VAL + row
		rowLabel := string(rune(rowLabelVal))
		rowData := []string{rowLabel}
		for col := range BOARD_SIZE {
			ship := board.sqaures[row][col]
			if ship == nil {
				rowData = append(rowData, OPTIONAL_SQUARE)
			} else {
				rowData = append(rowData, ship.icon)
			}
		}
		data = append(data, rowData)
	}
	return data
}

func (gs *gameState) ShowOpponentBoard(opponentBoard displayBoard) {
	defer fmt.Println()
	var data [][]string
	header := []string{" ", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	data = append(data, header)
	for row := range BOARD_SIZE {
		rowLabelVal := A_VAL + row
		rowLabel := string(rune(rowLabelVal))
		rowData := []string{rowLabel}
		for col := range BOARD_SIZE {
			square := opponentBoard.sqaures[row][col]
			if square == "" {
				square = OPTIONAL_SQUARE
			}
			rowData = append(rowData, square)
		}
		data = append(data, rowData)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Header(data[0])
	table.Bulk(data[1:])
	table.Render()

}

func (gs *gameState) pause() {
	fmt.Println("Game is paused.")
}

func NewGameState() *gameState {
	gs := &gameState{}
	return gs
}

var ErrInvalidOrientation = errors.New("error: not placed left to right or up and down")
var ErrInvalidShipSpace = errors.New("error: ship does not fit between start and end")

func validShipRange(sp shipPlacement) error {

	if sp.orientation == horizantal {
		if sp.start.col > sp.end.col {
			return ErrInvalidOrientation
		}
		if Abs(sp.start.col-sp.end.col)+1 != sp.ship.length {
			return ErrInvalidShipSpace
		}
		return nil
	} else {
		if sp.start.row > sp.end.row {
			return ErrInvalidOrientation
		}
		if Abs(sp.start.row-sp.end.row)+1 != sp.ship.length {
			return ErrInvalidShipSpace
		}
		return nil
	}
}

var ErrInvalidOccupiedSqaure = fmt.Errorf("error: there are ships already between start and end")

func (gs *gameState) shipsOccupyRange(sp shipPlacement, board board) error {
	if sp.orientation == horizantal {
		for i := 0; i < sp.ship.length; i++ {
			occupying := board.sqaures[sp.start.row][sp.start.col+i]
			if occupying != nil {
				return ErrInvalidOccupiedSqaure
			}
		}
		return nil
	} else {
		for i := 0; i < sp.ship.length; i++ {
			occupying := board.sqaures[sp.start.row+i][sp.start.col]
			if occupying != nil {
				return ErrInvalidOccupiedSqaure
			}
		}
		return nil
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func PrintServerHelp() {
	fmt.Println("help: prints possible commands")
	fmt.Println("log: prints server logs")
}

func PauseHandler(gs *gameState) func(msg routing.PauseMessage) {
	return func(msg routing.PauseMessage) {
		defer fmt.Print(">>> ")
		gs.pause()
		fmt.Println("Received pause message:", msg.Content)
	}
}

func NewPlayerHandler(gs *gameState) func(msg NewPlayerMessage) pubsub.AckType {
	return func(msg NewPlayerMessage) pubsub.AckType {
		fmt.Printf("New player joined: %s\n", msg.UserName)
		if gs.player1.userName == "" {
			gs.player1 = CreatePlayer(msg.UserName)
			gs.currentPlayer = gs.player1
			gs.player1Board.owner = gs.player1.userName
			fmt.Printf("Assigned %s as Player 1\n", msg.UserName)
			return pubsub.Ack
		} else if gs.player2.userName == "" {
			gs.player2 = CreatePlayer(msg.UserName)
			gs.player2Board.owner = gs.player2.userName
			fmt.Printf("Assigned %s as Player 2\n", msg.UserName)
			return pubsub.Ack
		} else {
			fmt.Printf("Both player slots are full. Could not assign %s\n", msg.UserName)
			return pubsub.NackDiscard
		}
	}
}

func ShowBoardHandler(gs *gameState, ch *amqp.Channel) func(msg routing.ShowBoardMessage) pubsub.AckType {
	return func(msg routing.ShowBoardMessage) pubsub.AckType {
		player := gs.getPlayerByName(msg.UserName)
		if player == nil {
			fmt.Printf("Could not find player with name: %s\n", msg.UserName)
			return pubsub.NackDiscard
		}
		var boardToShow board
		if gs.player1.userName == player.userName {
			boardToShow = gs.player1Board
		} else if gs.player2.userName == player.userName {
			boardToShow = gs.player2Board
		} else {
			fmt.Printf("Could not find board for player: %s\n", msg.UserName)
			return pubsub.NackDiscard
		}

		boardData := gs.boardData(boardToShow)

		err := pubsub.PublishJSON(ch, routing.EXCHANGE_BATTLESHIP_DIRECT, routing.BOARD_STATE_KEY, routing.ShowBoardMessage{UserName: msg.UserName, BoardData: boardData})
		if err != nil {
			fmt.Printf("Failed to publish board state message: %v\n", err)
			return pubsub.NackRequeue
		}
		fmt.Printf("Showing board for player: %s\n", msg.UserName)

		return pubsub.Ack
	}
}

func (gs *gameState) getPlayerByName(name string) *Player {
	if gs.player1.userName == name {
		return &gs.player1
	} else if gs.player2.userName == name {
		return &gs.player2
	}
	return nil
}

func ClientBoardStateHandler(userName string) func(msg routing.ShowBoardMessage) pubsub.AckType {
	return func(msg routing.ShowBoardMessage) pubsub.AckType {
		if msg.UserName != userName {
			return pubsub.NackDiscard
		}
		fmt.Println("Received board state update:")
		Show(msg.BoardData)
		return pubsub.Ack
	}
}
