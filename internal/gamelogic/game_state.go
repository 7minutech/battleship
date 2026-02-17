package gamelogic

const BOARD_SIZE = 10

const STARTING_SHIP_COUNT = 5

const CRUISER_LENGTH, SUBMARINE_LENGTH = 3, 3
const BATTLESHIP_LENGTH = 4
const CARRIER_LENGTH = 5
const DESTROYER_LENGTH = 2
const A_VAL = 97
const SHIP_CHAR = "●"
const OPTIONAL_SQUARE = "○"
const HIT_CHAR = "✖"
const MISS_CHAR = "~"
const START_ROW_HEADER = 1
const END_ROW_HEADER = 10
const START_COL_HEADER = 'a'
const END_COL_HEADER = 'j'

type PlayerTurn int

const (
	player1 PlayerTurn = 0
	player2 PlayerTurn = 1
)

type orientation int

const (
	horizontal orientation = 0
	vertical   orientation = 1
)

type gameState struct {
	player1       Player
	player2       Player
	currentPlayer Player
	player1Board  board
	player2Board  board
	gameOver      bool
}

type board struct {
	owner   string
	sqaures [BOARD_SIZE][BOARD_SIZE]*ship
}

type displayBoard struct {
	owner   string
	sqaures [BOARD_SIZE][BOARD_SIZE]string
}

type Player struct {
	userName  string
	shipCount int
	ships     []ship
}

type ship struct {
	name    string
	length  int
	start   boardMove
	end     boardMove
	modules map[boardMove]bool // a1[false], a2[false], a3[false], a4[true]
	hp      int
	alive   bool
	icon    string
}

type boardMove struct {
	row int
	col int
}

type shipPlacement struct {
	start       boardMove
	end         boardMove
	ship        ship
	orientation orientation
}

type NewPlayerMessage struct {
	UserName string `json:"user_name"`
}
