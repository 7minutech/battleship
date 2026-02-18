package routing

const EXCHANGE_BATTLESHIP_DIRECT = "battleship_direct"
const EXCHANGE_BATTLESHIP_TOPIC = "battleship_topic"
const PAUSE_KEY = "pause_started"
const NEW_PLAYER_KEY = "player.created"
const SHOW_BOARD_KEY = "board.show"
const GAME_COMMANDS_KEY = "game.commands"
const GAME_SERVER_KEY = "game.server"
const BOARD_STATE_KEY = "board.state"
const PLACE_STATE_KEY = "board.place.state"
const PLACE_BOARD_KEY = "board.place"
const AUTO_PLACE_KEY = "board.auto.place"
const GAME_RESET_KEY = "game.reset"
const AUTO_PLACE_STATE_KEY = "board.auto.place.state"
const OPPONENT_BOARD_STATE_KEY = "board.opponent.state"

type ShowBoardMessage struct {
	UserName  string
	BoardData [][]string
}

type PlaceShipCommand struct {
	UserName   string
	ShipType   string
	StartCoord string
	EndCoord   string
}

type AutoPlaceMessage struct {
	UserName string
	Success  bool
	Message  string
}

type ServerMessage struct {
	UserName string
	Success  bool
	Message  string
}

type ClientMessage struct {
	UserName string
}

type PlaceShipMessage struct {
	UserName string
	Success  bool
	Message  string
}

type PauseMessage struct {
	Content string
}

type GameResetMessage struct {
	Content string
}
