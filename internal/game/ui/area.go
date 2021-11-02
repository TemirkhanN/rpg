package ui

import "github.com/gdamore/tcell"

type Area struct {
	borderStyle tcell.Style
	rectangle   Box
}

func NewArea(x1 int, y1 int, x2 int, y2 int, borderStyle tcell.Style) Area {
	return Area{
		borderStyle: borderStyle,
		rectangle:   NewBox(x1, y1, x2, y2),
	}
}

func (a Area) Draw(on tcell.Screen) {
	a.rectangle.Draw(on, a.borderStyle)
}
