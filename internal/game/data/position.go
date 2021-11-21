package data

type Position interface {
	X() int
	Y() int
}

type Pos struct {
	x int
	y int
}

func (p Pos) X() int {
	return p.x
}

func (p Pos) Y() int {
	return p.y
}

func NewPos(x int, y int) Position {
	return Pos{
		x: x,
		y: y,
	}
}

var NoPosition = NewPos(0, 0)
