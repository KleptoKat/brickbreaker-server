package logic

import (
	"github.com/chrislonng/starx"
	"github.com/chrislonng/starx/timer"
	"time"
)

type gameService struct {
	gameIdIndex int64
	games map[int64]*Game
	gameTimer *timer.Timer
}

var gs = newGameService()

func newGameService() (gs *gameService) {
	gs = &gameService {
		0,
		make(map[int64]*Game, 0),
		timer.Register(time.Millisecond * 10,updateGames),
	}
	return
}


func getNextGameID() (id int64) {
	id = gs.gameIdIndex
	gs.gameIdIndex += 1
	return
}

func RegisterGame(p1 int64, p2 int64) (game *Game) {

	game = NewGame(p1, p2)
	starx.Register(game)
	gs.games[game.ID] = game
	return
}


func NewGame(p1 int64, p2 int64) (game *Game) {

	id := getNextGameID()

	game = &Game{
		ID:      id,
		P1ID:    p1,
		P2ID:    p2,
		channel: starx.ChannelService.NewChannel("Game_" + string(id)),
	}

	return
}

func EndGameByID(id int64) {
	FindGameByID(id).End()
}

func GetGames() (games []*Game) {
	games = make([]*Game, len(gs.games))
	return
}


func FindGameByID(id int64) *Game {
	return gs.games[id]
}

func CountGames() int {
	return len(GetGames())
}

func updateGames() {
	for _, game := range GetGames() {
		game.update()
	}
}

func StopUpdating() {

	gs.gameTimer.Stop()
}