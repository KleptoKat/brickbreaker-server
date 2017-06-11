package main

import (
	"net/http"

	"github.com/chrislonng/starx"
	"github.com/chrislonng/starx/log"
	"github.com/chrislonng/starx/serialize/json"
	"github.com/KleptoKat/brickbreaker-server/account"
	"github.com/KleptoKat/brickbreaker-server/logic"
	"github.com/KleptoKat/brickbreaker-server/matchmaking"
)

func main() {
	starx.SetAppConfig("configs/app.json")
	starx.SetServersConfig("configs/servers.json")
	starx.Register(account.NewManager())
	starx.Register(matchmaking.NewMatchmaker())

	starx.SetServerID("brickbreaker-server-1")
	starx.SetSerializer(json.NewSerializer())

	log.SetLevel(log.LevelDebug)

	starx.SetCheckOriginFunc(func(_ *http.Request) bool { return true })

	starx.Run()
	logic.StopUpdating()
}
