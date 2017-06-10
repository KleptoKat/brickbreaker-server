package logic

import (
	"github.com/chrislonng/starx"
	"errors"
	"github.com/chrislonng/starx/timer"
	"time"
)

var games []*Game
var gameTimer *timer.Timer
var gameCount int64 = 0

func RegisterGame(p1 int64, p2 int64) (game *Game) {

	game = NewGame(p1, p2)
	starx.Register(game)
	return
}


func NewGame(p1 int64, p2 int64) (game *Game) {

	game = &Game{
		ID:       gameCount,
		P1ID:     p1,
		P2ID:     p2,
		channel:  starx.ChannelService.NewChannel("Game"),
	}

	gameCount++

	return
}

func EndGame(id int64) {
	/*game, index, err := findGameByID(id)
	game.End()
	games = append(games[:index], games[index+1:]...)*/
}

func findGameByID(id int64) (*Game, int,  error) {

	for index, game := range games {
		if game.ID == id {
			return game, index, nil
		}
	}

	return nil, -1, errors.New("Could not find game with id " + string(id))
}

func CountGames() int {
	return len(games)
}

func StartUpdating() error {
	/*for _, game := range games {
		// update game here
	}*/

	gameTimer = timer.Register(time.Millisecond * 10, func() {
		for _, game := range games {
			game.update()
		}
	})
	return nil
}

func StopUpdating() {

	gameTimer.Stop()
}