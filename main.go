package main

import (
	"net/http"

	"github.com/chrislonng/starx"
	"github.com/chrislonng/starx/log"
	"github.com/chrislonng/starx/serialize/json"
	"github.com/KleptoKat/brickbreaker-server/account"
	"github.com/KleptoKat/brickbreaker-server/logic"
	"github.com/KleptoKat/brickbreaker-server/matchmaking"
	"github.com/KleptoKat/brickbreaker-server/database"
	"time"
	"github.com/chrislonng/starx/session"
)

var mm *matchmaking.Matchmaker
var gs *logic.GameService

func main() {
	starx.SetAppConfig("configs/app.json")
	starx.SetServersConfig("configs/servers.json")
	starx.Register(account.NewManager())

	mm = matchmaking.NewMatchmaker()
	gs = logic.NewGameService()

	starx.Register(mm)
	starx.Register(gs)

	starx.SetServerID("brickbreaker-server-1")
	starx.SetSerializer(json.NewSerializer())

	log.SetLevel(log.LevelDebug)

	starx.SetCheckOriginFunc(func(_ *http.Request) bool { return true })
	starx.OnSessionClosed(OnSessionClosedCallback)

	database.OpenConnection()
	starx.SetHeartbeatInternal(2*time.Second)
	starx.Run()
	database.CloseConnection()
}


func OnSessionClosedCallback(s *session.Session) {
	mm.Remove(s)
	gs.Remove(s)


}