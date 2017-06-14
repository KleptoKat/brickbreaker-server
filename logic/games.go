package logic

import (
	"github.com/chrislonng/starx"
	"github.com/chrislonng/starx/timer"
	"time"
	"github.com/chrislonng/starx/component"
)

type GameService struct {


	component.Base
	channel *starx.Channel

	gameIdIndex int64
	games map[int64]*Game
	gameTimer *timer.Timer
}



type JoinRequest struct {
	GameID int64 `json:"game_id"`
}

var gs *GameService

func NewGameService() (*GameService) {
	gs = &GameService {
		games:make(map[int64]*Game, 0),
		gameIdIndex:0,
		gameTimer:nil,
		channel:starx.ChannelService.NewChannel("GameService"),
	}
	StartTimer()
	return gs
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
		State:"waiting",
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

func StartTimer() {

	gs.gameTimer = timer.Register(time.Millisecond * 15,updateGames)
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

func (gs *GameService) findGameByPlayerID(id int64) (*Game) {
	for _, game := range gs.games {
		if game.P1ID == id || game.P2ID == id {
			return game
		}
	}

	return nil
}

func (gs *GameService) findGameByInGamePlayerID(id int64) (*Game) {
	for _, game := range gs.games {
		if  (game.p1 != nil && game.p1.Uid == id) ||
			(game.p2 != nil && game.p2.Uid == id) {
			return game
		}
	}

	return nil
}