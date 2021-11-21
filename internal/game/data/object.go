package data

import "github.com/gdamore/tcell"

type Object struct {
	solid bool
	pos   Position
	icon  rune
}

func (o Object) Solid() bool {
	return o.solid
}

func (o Object) Collides(with Position) bool {
	if !o.Solid() {
		return false
	}

	return collides(o.Position(), with)
}

func (o Object) Position() Position {
	return o.pos
}

func (o Object) Icon() rune {
	return o.icon
}

func (o Object) Style() tcell.Style {
	return objectStyle
}

var objectStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorLightGray).Bold(true)
