package logic

import (
	"github.com/chrislonng/starx/component"
	"github.com/chrislonng/starx"
	"github.com/chrislonng/starx/session"
	"github.com/KleptoKat/brickbreaker-server/response"
)


type Game struct {
	ID int64 `json:"id"`
	P1ID int64 `json:"p1"`
	P2ID int64 `json:"p2"`

	p1 *session.Session
	p2 *session.Session

	component.Base
	channel *starx.Channel
}

type JoinRequest struct {
	GameID int64 `json:"game_id"`
}

type LeaveRequest struct {

}


func (game *Game) update() error {
	return nil
}

func (game *Game) End() error {
	game.Leave(game.p1, &LeaveRequest{})
	game.Leave(game.p2, &LeaveRequest{})

	return nil
}

func (game *Game) Join(s *session.Session, msg *JoinRequest) error {

	if s.Uid == game.P1ID {
		game.p1 = s
	} else if s.Uid == game.P2ID {
		game.p2 = s
	} else {
		return s.Response(response.BadRequest())
	}


	game.channel.Add(s) // add session to channel

	if game.p1 != nil && game.p2 != nil {

	}

	return nil
}

func (game *Game) Leave(s *session.Session, msg *LeaveRequest) error {
	if game.channel.IsContain(s.Uid) {
		game.End()
		return s.Response(response.OK())
	}

	return s.Response(response.BadRequest())
}
