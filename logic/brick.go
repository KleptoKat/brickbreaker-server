package logic

import "github.com/KleptoKat/brickbreaker-server/collision2d"

type BrickMessage struct {
	collision2d.Box `json:"box"`
	collision2d.Vector `json:"vel"`
}

type Brick struct {
	*BoxCollider
}

func (brick Brick) GetType() int {
	return BODY_TYPE_BRICK
}

func (brick Brick) GetJSONInterface() interface{} {
	return BrickMessage {
		brick.Box,
		brick.vel,
	}
}

func (brick Brick) GetCollisionBody() CollisionBody {
	return brick.BoxCollider
}

func NewBrick(pos collision2d.Vector, width float64, height float64) *Brick {
	return &Brick {
		BoxCollider: &BoxCollider{
			vel:collision2d.NewVector(0,0),
			Box:collision2d.NewBox(pos, width, height),
		},
	}
}