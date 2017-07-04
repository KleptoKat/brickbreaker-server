package matchmaking

import (
	"github.com/chrislonng/starx/component"
	"github.com/chrislonng/starx"
	"github.com/chrislonng/starx/session"
	"github.com/KleptoKat/brickbreaker-server/response"
	"github.com/chrislonng/starx/timer"
	"time"
	"github.com/KleptoKat/brickbreaker-server/logic"
	"github.com/chrislonng/starx/log"
)

type Matchmaker struct {

	queue []*session.Session
	maxGames int
	timer *timer.Timer

	component.Base
	channel *starx.Channel
}

type SearchStatus struct {
	Code int `json:"code"`
}

type SearchRequest struct {

}

type FoundGameMessage struct {
	GameID int64 `json:"game_id"`
}


func NewMatchmaker() (mm *Matchmaker) {
	mm = &Matchmaker{
		queue:    make([]*session.Session, 0),
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

	active := int64(0)

	for _, Uid := range mm.channel.Members() {

		if active == 0 {
			active = Uid
		} else if logic.CountGames() < mm.maxGames && active != Uid {
			mm.match(active, Uid)
			active = 0
		} else {
			break // breaks if there is players waiting and too many games playing.
		}
	}
}


func (mm *Matchmaker) match(p1 int64, p2 int64) {


	s1 := mm.channel.Member(p1)
	s2 := mm.channel.Member(p2)
	log.Debug("Matching " + string(p1) + " with " + string(p2))

	game := logic.RegisterGame(s1, s2)

	if s1 == nil || s2 == nil {
		return
	}

	startGame := true
	if err1 := sendGameToMember(s1, s2, game); err1 != nil {
		mm.Remove(s1)
		startGame = false
	}


	if err2 := sendGameToMember(s2, s1, game); err2 != nil {
		mm.Remove(s2)
		startGame = false
	}

	if startGame {
		mm.Remove(s1)
		mm.Remove(s2)
		log.Info("Match successful between " + string(s1.Uid) + " and " + string(s2.Uid))
		game.SetReady()
	}

}


func sendGameToMember(s *session.Session, opposingS *session.Session, g *logic.Game) error {

	msg := FoundGameMessage {
		GameID:g.ID,
	}

	return s.Push("PrepGame", msg)
}


func (mm *Matchmaker) getSearchStatus(s *session.Session) SearchStatus {

	if s.Uid == 0 {
		return SearchStatus{ Code:0 } // not logged in
	} else if mm.channel.IsContain(s.Uid) {
		return SearchStatus{ Code:1  } // searching
	} else {
		return SearchStatus{ Code:2 } // not searching
	}
}


func (mm *Matchmaker) StartSearching(s *session.Session, msg *SearchRequest) error {


	if s.Uid == 0 {
		return s.Response(response.BadRequestWithDescription("Not logged in"))
	}

	if mm.channel.IsContain(s.Uid) {
		return s.Response(response.BadRequestWithDescription("Already searching with ID " + string(s.Uid)))
	}

	mm.channel.Add(s)
	mm.queue = append(mm.queue, s)


	return s.Response(response.OKWithData(mm.getSearchStatus(s)))
}


func (mm *Matchmaker) Remove(s *session.Session)  {
	if mm.channel.IsContain(s.Uid) {
		mm.channel.Leave(s.Uid)

		// remove from queue
		for i,mem := range mm.queue {
			if mem.Uid == s.Uid {
				mm.queue = append(mm.queue[:i], mm.queue[i+1:]...)
				break
			}
		}
	}
}

func (mm *Matchmaker) StopSearching(s *session.Session, msg *SearchRequest) error {
	mm.Remove(s)
	return nil
}


func (mm *Matchmaker) RetrieveSearchStatus(s *session.Session, msg *SearchRequest) error {

	if s.Uid == 0 {
		return s.Response(response.BadRequest())
	}

	return s.Response(response.OKWithData(mm.getSearchStatus(s)))

}