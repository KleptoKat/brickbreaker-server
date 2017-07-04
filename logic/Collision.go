package logic


import "github.com/KleptoKat/brickbreaker-server/collision2d"

type CollisionBody interface {
	SetVelocity(collision2d.Vector)
	GetVelocity() collision2d.Vector
	CollidesWithCircle(other CircleCollider) (bool, collision2d.Response)
	CollidesWithBox(other BoxCollider) (bool, collision2d.Response)
}


//////////////////////////////
//////////// BOX /////////////
//////////////////////////////

type BoxCollider struct {
	vel collision2d.Vector `json:"vel"`
	collision2d.Box `json:"box"`
}

func (col *BoxCollider) GetVelocity() collision2d.Vector {
	return col.vel
}

func (col *BoxCollider) SetVelocity(vel collision2d.Vector) {
	col.vel = vel
}

func (col *BoxCollider) CollidesWithCircle(other CircleCollider) (bool, collision2d.Response) {
	return collision2d.TestPolygonCircle(col.Box.ToPolygon(), other.Circle)
}

func (col *BoxCollider) CollidesWithBox(other BoxCollider) (bool, collision2d.Response) {
	return collision2d.TestPolygonPolygon(col.Box.ToPolygon(), other.Box.ToPolygon())
}


//////////////////////////////
/////////// CIRCLE ///////////
//////////////////////////////

type CircleCollider struct {
	vel collision2d.Vector
	collision2d.Circle
}


func (col CircleCollider) GetVelocity() collision2d.Vector {
	return col.vel
}

func (col *CircleCollider) SetVelocity(vel collision2d.Vector) {
	col.vel = vel
}

func (col CircleCollider) CollidesWithCircle(other CircleCollider) (bool, collision2d.Response) {
	return collision2d.TestCircleCircle(col.Circle, other.Circle)
}

func (col CircleCollider) CollidesWithBox(other BoxCollider) (bool, collision2d.Response) {
	return collision2d.TestCirclePolygon(col.Circle, other.ToPolygon())
}
