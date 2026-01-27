package gamelogic

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func GetWords(input string) []string {
	words := strings.Split(input, " ")
	return words
}

func convertMove(mv string) (boardMove, error) {
	if len(mv) != 2 {
		return boardMove{}, fmt.Errorf("error: move was not exactly 2: %s", mv)
	}

	letter := mv[0]
	num := string(mv[1])

	numInt, err := strconv.Atoi(num)
	if err != nil {
		return boardMove{}, fmt.Errorf("error: could not convert number in move to int: %s", mv)
	}

	if !(int('a') <= int(letter) && int(letter) <= int('z')) {
		return boardMove{}, fmt.Errorf("error: move did not contain a letter that was between a-z: %s", mv)
	}

	if !(1 <= numInt && numInt <= 9) {
		return boardMove{}, fmt.Errorf("error: move did not contain a number that was between 1-9: %s", mv)
	}

	row := (int(letter) - int('a'))
	col := numInt - 1

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

func (gs *gameState) PlaceShip(words []string) error {
	sp, err := gs.getShipPlacement(words)
	if err != nil {
		return err
	}
	if err := gs.validateShipPlacement(sp); err != nil {
		return err
	}

	if sp.orientation == horizantal {
		for i := range sp.ship.length {
			row := sp.start.row
			col := sp.start.col + i
			gs.gameBoard.sqaures[row][col] = &sp.ship
		}
	} else {
		for i := range sp.ship.length {
			row := sp.start.row + i
			col := sp.start.col
			gs.gameBoard.sqaures[row][col] = &sp.ship
		}
	}

	return nil
}

func (gs *gameState) getShipPlacement(words []string) (shipPlacement, error) {
	shipName := words[1]
	startPlace := words[2]
	endPlace := words[3]
	ship, err := gs.getShip(shipName)
	if err != nil {
		return shipPlacement{}, err
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

func (gs *gameState) validateShipPlacement(sp shipPlacement) error {
	var err error
	err = validShipRange(sp)
	err = gs.shipsOccupyRange(sp)
	return err
}

func (gs *gameState) getShip(shipName string) (ship, error) {
	for _, ship := range gs.player.ships {
		if ship.name == shipName {
			return ship, nil
		}
	}
	return ship{}, fmt.Errorf("error: could not find ship: %s", shipName)
}

func (gs *gameState) Show() {
	defer fmt.Println()
	var data [][]string
	header := []string{" ", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	data = append(data, header)
	for row := range BOARD_SIZE {
		rowLabelVal := A_VAL + row
		rowLabel := string(rune(rowLabelVal))
		rowData := []string{rowLabel}
		for col := range BOARD_SIZE {
			ship := gs.gameBoard.sqaures[row][col]
			if ship == nil {
				rowData = append(rowData, OPTIONAL_SQUARE)
			} else {
				rowData = append(rowData, ship.icon)
			}
		}
		data = append(data, rowData)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.Header(data[0])
	table.Bulk(data[1:])
	table.Render()

}

func NewGameState(player Player) *gameState {
	var gs gameState = gameState{player: player, turn: player1}
	return &gs
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

func (gs *gameState) shipsOccupyRange(sp shipPlacement) error {
	if sp.orientation == horizantal {
		for i := 0; i < sp.ship.length; i++ {
			occupying := gs.gameBoard.sqaures[sp.start.row+i][sp.start.col]
			if occupying != nil {
				return ErrInvalidOccupiedSqaure
			}
		}
		return nil
	} else {
		for i := 0; i < sp.ship.length; i++ {
			occupying := gs.gameBoard.sqaures[sp.start.row][sp.start.col+i]
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
