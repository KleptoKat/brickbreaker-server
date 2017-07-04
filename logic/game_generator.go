package logic

import "github.com/KleptoKat/brickbreaker-server/collision2d"

const (
	DEFAULT_GAME_BRICK_COUNT_X  = 6
	DEFAULT_GAME_BRICK_COUNT_Y  = 5
	DEFAULT_GAME_BRICK_WIDTH  = 500
	DEFAULT_GAME_BRICK_HEIGHT   = 200
	DEFAULT_GAME_BRICK_SPACING  = 500


	DEFAULT_GAME_PADDLE_HEIGHT  = 50
	DEFAULT_GAME_PADDLE_WIDTH   = 100
	DEFAULT_GAME_PADDLE_ALT     = 60


)

func generateDefaultGame() (game *Game) {

	game = &Game {
		width:GAME_WORLD_WIDTH,
		height:GAME_WORLD_HEIGHT,
		bodies:make(map[int]*Body, 0),
		bodyIdIndex:1,
	}

	brick := *NewBrick(
		collision2d.NewVector(
			float64(game.width/2 - DEFAULT_GAME_BRICK_WIDTH/2),
			float64(game.height/2 - DEFAULT_GAME_BRICK_HEIGHT/2)),
		DEFAULT_GAME_BRICK_WIDTH, DEFAULT_GAME_BRICK_HEIGHT)

	game.addBodyNow(brick)


	// do game generation here.


	return
}