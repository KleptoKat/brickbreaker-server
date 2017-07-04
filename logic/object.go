package logic

const (
	BODY_TYPE_BRICK            = 0
	BODY_TYPE_BALL             = 1
	BODY_TYPE_PADDLE           = 2
	BODY_TYPE_PLAYER_PADDLE    = 3
)


type Body interface {
	GetCollisionBody() CollisionBody
	GetJSONInterface() interface{}
	GetType() int
}

type Ball struct {
	*CircleCollider
}

func (ball *Ball) GetType() int {
	return BODY_TYPE_BALL
}

func (ball *Ball) GetJSONInterface() interface{} {
	panic("implement me")
}

func (ball *Ball) GetCollisionBody() CollisionBody {
	return ball.CircleCollider
}

type Paddle struct {
	*BoxCollider
}

func (paddle *Paddle) GetType() int {
	return BODY_TYPE_PADDLE
}

func (paddle *Paddle) GetJSONInterface() interface{} {
	panic("implement me")
}

func (paddle *Paddle) GetCollisionBody() CollisionBody {
	return paddle.BoxCollider
}



