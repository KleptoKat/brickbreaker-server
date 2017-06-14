package logic

import (
	"github.com/chrislonng/starx/component"
	"github.com/chrislonng/starx"
	"github.com/chrislonng/starx/session"
	"github.com/KleptoKat/brickbreaker-server/response"
	"github.com/KleptoKat/brickbreaker-server/account"
)


type Game struct {
	ID int64 `json:"id"`
	P1ID int64 `json:"p1"`
	P2ID int64 `json:"p2"`
	State string `json:"state"`

	p1 *session.Session
	p2 *session.Session

	component.Base
	channel *starx.Channel
}

type GameMessage struct {
	GameID int64 `json:"game_id"`
	OpposingName string `json:"opposing_name"`
	State string `json:"state"`
}


type GameState struct {
	State string `json:"state"`
}

type LeaveRequest struct {}


func (game *Game) update() {
	/*if game.p1 != nil && game.p2 != nil {
		game.State = "ready"
		game.onStateChange()
	}*/
}

func (game *Game) End() {
	game.Leave(game.p1, &LeaveRequest{})
	game.Leave(game.p2, &LeaveRequest{})
	game.channel.Destroy()
	delete(gs.games, game.ID)
}

func (game *Game) onStateChange() {
	game.channel.Broadcast("GameStateMessage", &GameState {
		State:game.State,
	})
}

func (game *Game) Leave(s *session.Session, msg *LeaveRequest) error {
	if game == nil {
		return s.Response(response.BadRequestWithDescription("No game"))

	}

	s.State()

	if game.channel.IsContain(s.Uid) {
		game.End()
		return s.Response(response.OK())
	}

	return s.Response(response.BadRequest())
}

func (gs *GameService) JoinGame(s *session.Session, msg *JoinRequest) error {

	game := FindGameByID(msg.GameID)
	if game == nil {
		return s.Response(response.BadRequestWithDescription("No game"))
	}

	if gs.findGameByInGamePlayerID(s.Uid) != nil {
		return s.Response(response.BadRequestWithDescription("Already in game."))
	}

	var opposingID int64;
	if s.Uid == game.P1ID {
		game.p1 = s
		opposingID = game.P2ID;
	} else if s.Uid == game.P2ID {
		game.p2 = s
		opposingID = game.P1ID;
	} else {
		return s.Response(response.BadRequest())
	}


	game.channel.Add(s) // add session to channel

	return s.Response(response.OKWithData(&GameMessage{
		GameID:game.ID,
		OpposingName:*account.GetNameByID(opposingID),
		State:game.State,
	}))
}

func (gs *GameService) Leave(s *session.Session, msg *LeaveRequest) error {

	game := gs.findGameByPlayerID(s.Uid)
	if game == nil {
		return s.Response(response.BadRequestWithDescription("No game"))

	}

	s.State()

	if game.channel.IsContain(s.Uid) {
		game.End()
		return s.Response(response.OK())
	}

	return s.Response(response.BadRequest())
}