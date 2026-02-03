package routing

const EXCHANGE_BATTLESHIP_DIRECT = "battleship_direct"
const PAUSE_KEY = "pause_started"
const NEW_PLAYER_KEY = "player.created"
const SHOW_BOARD_KEY = "board.show"
const GAME_COMMANDS_QUEUE = "game.commands"
const GAME_SERVER_KEY = "game.server"
const BOARD_STATE_KEY = "board.state"

type ShowBoardMessage struct {
	UserName  string
	BoardData [][]string
}

type PauseMessage struct {
	Content string
}
