package ui

import "github.com/gdamore/tcell"

type SolidObject struct {
	asci asci
	pos  Position
}

func NewSolidObject(icon rune, iconStyle tcell.Style, pos Position) SolidObject {
	return SolidObject{
		asci: asci{style: iconStyle, symbol: icon},
		pos:  pos,
	}
}

func (s SolidObject) Draw(on tcell.Screen) {
	NewText(string(s.asci.symbol), s.pos.X, s.pos.Y).Draw(on, s.asci.style)
}

func (s SolidObject) Collides(with Position) bool {
	if with.Y != s.pos.Y {
		return false
	}

	if with.X < s.pos.X-1 || with.X > s.pos.X+1 {
		return false
	}

	return true
}
