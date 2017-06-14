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
)

func main() {
	starx.SetAppConfig("configs/app.json")
	starx.SetServersConfig("configs/servers.json")
	starx.Register(account.NewManager())
	starx.Register(matchmaking.NewMatchmaker())
	starx.Register(logic.NewGameService())

	starx.SetServerID("brickbreaker-server-1")
	starx.SetSerializer(json.NewSerializer())

	log.SetLevel(log.LevelDebug)

	starx.SetCheckOriginFunc(func(_ *http.Request) bool { return true })
	database.OpenConnection()
	starx.SetHeartbeatInternal(5*time.Second)
	starx.Run()
	database.CloseConnection()
}
