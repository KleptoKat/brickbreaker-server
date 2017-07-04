package logic

import (
	"github.com/chrislonng/starx"
	"github.com/chrislonng/starx/timer"
	"time"
	"github.com/chrislonng/starx/component"
	"github.com/chrislonng/starx/session"
	"github.com/KleptoKat/brickbreaker-server/response"
	"github.com/KleptoKat/brickbreaker-server/account"
	"github.com/chrislonng/starx/log"
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
		gameIdIndex:1,
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

func RegisterGame(p1 *session.Session, p2 *session.Session) (game *Game) {

	game = NewGame(p1, p2)
	gs.games[game.ID] = game
	return
}


func NewGame(p1 *session.Session, p2 *session.Session) (game *Game) {

	game = generateDefaultGame()
	game.ID = getNextGameID()
	game.channel = starx.ChannelService.NewChannel("Game_" + string(game.ID))
	game.p1 = p1
	game.p2 = p2
	game.State = GAME_STATE_PREP
	game.lastState = game.State
	game.createTime = time.Now()

	return
}

func EndGameByID(id int64) {
	FindGameByID(id).End()
}

func GetGames() (map[int64]*Game) {
	return gs.games
}

func StartTimer() {

	gs.gameTimer = timer.Register(time.Millisecond * 15,updateGames)
}

func FindGameByID(id int64) *Game {
	for id, game := range gs.games {
		log.Debug("ID: " + string(id))
		if game == nil {
			log.Debug("GAME NIL")
		}
	}
	return gs.games[id]
}

func CountGames() int {
	return len(GetGames())
}

func updateGames() {
	for _, game := range GetGames() {
		if game != nil {
			game.update()
		}
	}
}

func (gs *GameService) findGameByPlayerID(id int64) (*Game) {
	for _, game := range gs.games {
		if game.p1.Uid == id || game.p2.Uid == id {
			return game
		}
	}

	return nil
}

func (gs *GameService) Remove(s *session.Session) {
	game := gs.findGameByPlayerID(s.Uid)
	if game != nil && s != nil {
		game.onPlayerDisconnect(s)
	}
}



func (game *Game) onPlayerDisconnect(s *session.Session) {
	//log.Info("Player " + string(s.Uid) + " disconnected from game " + string(game.ID))
	log.Info("Player disconnected. " + string(game.State))
	if game.State != GAME_STATE_ENDED {
		game.End()
	}
}

func (game *Game) Leave(s *session.Session, msg *LeaveRequest) error {
	if game == nil {
		return s.Response(response.BadRequestWithDescription("No game"))

	}

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

	var opposingID int64
	switch s.ID {
	case game.p1.ID:
		opposingID = game.p2.Uid
		break
	case game.p2.ID:
		opposingID = game.p1.Uid
		break
	default:
		return s.Response(response.BadRequest())
		break
	}

	game.channel.Add(s) // add session to channel

	return s.Response(response.OKWithData(&FullGameUpdate{
		GameID:game.ID,
		OpposingName:account.GetNameByID(opposingID),
		State:game.State,
		Width:game.width,
		Height:game.height,
		Bodies:game.getBodiesWithType(),
	}))
}

func (gs *GameService) Leave(s *session.Session, msg *LeaveRequest) error {

	game := gs.findGameByPlayerID(s.Uid)
	if game == nil {
		return nil
	}

	if game.channel.IsContain(s.Uid) {
		game.End()
	}

	return nil
}