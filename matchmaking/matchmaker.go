package matchmaking

import (
	"github.com/chrislonng/starx/component"
	"github.com/chrislonng/starx"
	"github.com/chrislonng/starx/session"
	"github.com/KleptoKat/brickbreaker-server/response"
	"github.com/chrislonng/starx/timer"
	"time"
	"github.com/KleptoKat/brickbreaker-server/logic"
)

type Matchmaker struct {

	queue []session.Session
	maxGames int
	timer *timer.Timer

	component.Base
	channel *starx.Channel
}

type SearchStatus struct {
	Code int `json:"code"`
	Status string `json:"status"`
}

type SearchRequest struct {

}

type FoundGameMessage struct {
	GameID int64 `json:"game_id"`
}


func NewMatchmaker() (mm *Matchmaker) {
	mm = &Matchmaker{
		queue:    make([]session.Session, 0),
		maxGames: 50,
		channel:  starx.ChannelService.NewChannel("Matchmaker"),
		timer:    nil,
	}
	defer mm.startUpdateTimer()
	return mm
}

func (mm *Matchmaker) startUpdateTimer() {
	mm.timer = timer.Register(time.Second, mm.update)
}

func (mm *Matchmaker) update() {

	active := int64(0);
	for _, Uid := range mm.channel.Members() {
		if active == 0 {
			active = Uid
		} else if logic.CountGames() < mm.maxGames && active != Uid {
			mm.match(active, Uid)
			active = 0
		} else {
			break // breaks if there is players waiting and
		}
	}

}


func (mm *Matchmaker) match(p1 int64, p2 int64) {

	game := logic.RegisterGame(p1, p2)

	s1 := mm.channel.Member(p1)
	s2 := mm.channel.Member(p2)

	if s1 == nil || s2 == nil {
		return
	}

	sendGameToMember(s1, s2, game)
	sendGameToMember(s2, s1, game)
}


func sendGameToMember(s *session.Session, opposingS *session.Session, g *logic.Game) error {

	msg := FoundGameMessage {
		GameID:g.ID,
	}

	return s.Push("JoinGame", msg)
}


func (mm *Matchmaker) getSearchStatus(s *session.Session) SearchStatus {

	if s.Uid == 0 {
		return SearchStatus{ Code:0,Status:"Not logged in", }
	} else if mm.channel.IsContain(s.Uid) {
		return SearchStatus{ Code:1,Status:"Searching...", }
	} else {
		return SearchStatus{ Code:2,Status:"Not searching", }
	}
}


func (mm *Matchmaker) StartSearching(s *session.Session, msg *SearchRequest) error {
	return s.Response(response.BadRequest())
}

func (mm *Matchmaker) RetrieveSearchStatus(s *session.Session, msg *SearchRequest) error {

	if s.Uid == 0 {
		return s.Response(response.BadRequest())
	}

	return s.Response(response.OKWithData(mm.getSearchStatus(s)))

}