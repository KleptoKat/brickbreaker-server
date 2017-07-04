package logic

import (
	"github.com/chrislonng/starx/component"
	"github.com/chrislonng/starx"
	"github.com/chrislonng/starx/session"
	"time"
)

const (
	GAME_STATE_PREP          = 0
	GAME_STATE_READY         = 1
	GAME_STATE_STARTING      = 2
	GAME_STATE_RUNNING       = 3
	GAME_STATE_ENDED         = 4
	GAME_STATE_CANCELED      = 5

	GAME_START_INTERVAL      = 3000
	GAME_PREP_TIMEOUT        = 6000

	GAME_WORLD_WIDTH         = 900
	GAME_WORLD_HEIGHT        = 1600
)

type Game struct {
	ID int64 `json:"id"`
	State int `json:"state"`
	lastState int
	ready bool

	createTime time.Time
	startTime time.Time

	lastUpdate time.Time
	lastNetworkUpdate time.Time

	p1 *session.Session
	p2 *session.Session

	height int
	width int

	bodies map[int]*Body
	bodyIdIndex int

	component.Base
	channel *starx.Channel
}

type FullGameUpdate struct {
	GameID int64 `json:"game_id"`
	OpposingName string `json:"opposing_name"`
	State int `json:"state"`
	Width int `json:"width"`
	Height int `json:"height"`
	Bodies []bodyWithTypeAndID `json:"bodies"`
}


type GameState struct {
	State int `json:"state"`
}

type LeaveRequest struct {}


type StartTimeMessage struct {
	Start int64 `json:"start"`
	Duration int `json:"dur"` // in milliseconds

}

func (game *Game) bodyUpdate() {
	//delta := time.Now().Sub(game.lastUpdate)

}

func (game *Game) networkUpdate() {
	//delta := time.Now().Sub(game.lastUpdate)

}

func (game *Game) update() {

	// do body update
	// do network update


	game.lastState = game.State
	game.updateState()
	if game.lastState != game.State {
		game.onStateChange()
	}
}

func (game *Game) updateState() {
	switch game.State {
	case GAME_STATE_PREP:
		// Will set game state to ready when all players have joined and the matchmaker is ready.
		if game.channel.IsContain(game.p1.Uid) && game.channel.IsContain(game.p2.Uid) && game.ready {
			game.State = GAME_STATE_READY
			break
		}

		// If prep lasts for too long the game will time out.
		if game.createTime.UnixNano() + GAME_PREP_TIMEOUT*1000000 < time.Now().UnixNano() {
			game.Cancel()
		}
		break
	case GAME_STATE_STARTING:
		if game.startTime.UnixNano() < time.Now().UnixNano() {
			game.run()
		}
		break
	}

}
func (game *Game) run() {

	game.State = GAME_STATE_RUNNING

}


func (game *Game) SetReady() {
	game.ready = true
}

func (game *Game) startIn(interval time.Duration) {
	game.startTime = time.Now().Add(interval)

	game.State = GAME_STATE_STARTING
	game.channel.Broadcast("GameStartTime", &StartTimeMessage {
		Start:game.startTime.UnixNano(),
		Duration:int(interval.Nanoseconds()/1000000),
	})
}

func (game *Game) onStateChange() {
	game.channel.Broadcast("GameStateMessage", &GameState {
		State:game.State,
	})

	switch game.State {
	case GAME_STATE_READY:
		game.startIn(GAME_START_INTERVAL * 1000000)
	}
}

func (game *Game) incrementBodyIndex() {
	game.bodyIdIndex = game.bodyIdIndex + 1
}

func (game *Game) getNextBodyID() int {
	defer game.incrementBodyIndex()
	return game.bodyIdIndex
}

/*
func (game *Game) getBricks() []*Brick {
	return game.bodies
}

func (game *Game) getP1Paddle() *Paddle {
	return game.bodies
}

func (game *Game) getP2Paddle() *Paddle {
	return game.bodies
}


func (game *Game) getBalls() []*Ball {
	return game.bodies
}*/

func (game *Game) getBodies() map[int]*Body {
	return game.bodies
}

type bodyWithTypeAndID struct {
	B *Body `json:"b"`
	ID int `json:"id"`
	T int `json:"t"`
}

func (game *Game) getBodiesWithType() []bodyWithTypeAndID {
	bodies := make([]bodyWithTypeAndID,0)

	var newBody bodyWithTypeAndID
	for id, body := range game.bodies {

		newBody = bodyWithTypeAndID{
			body,
			id,
			(*body).GetType(),
		}

		bodies = append(bodies, newBody)
	}

	return bodies
}



func (game *Game) addBodyNow(body Body) {
	game.bodies[game.getNextBodyID()] = &body
}

func (game *Game) remove() {
	game.channel.Destroy()
	delete(gs.games, game.ID)
}

func (game *Game) Cancel() {
	game.State = GAME_STATE_CANCELED
	game.onStateChange()
	game.remove()
}

func (game *Game) End() {
	game.State = GAME_STATE_ENDED
	game.onStateChange()
	game.remove()
}