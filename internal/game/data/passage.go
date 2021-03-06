package data

import (
	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/internal/game/ui/icon"
)

type Passage struct {
	in  Position
	out Position
	to  Location
}

func (p Passage) Position() Position {
	return p.in
}

func (p Passage) Icon() rune {
	return icon.Door
}

func (p Passage) Style() tcell.Style {
	return passageStyle
}

var passageStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBrown).Bold(true)
