package data

import (
	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type Npc struct {
	template rpg.NPC
	icon     rune
	position Position
}

func NewNpc(template rpg.NPC, icon rune, pos Position) Npc {
	return Npc{
		icon:     icon,
		template: template,
		position: pos,
	}
}

func (n Npc) Position() Position {
	return n.position
}

func (n Npc) Icon() rune {
	return n.icon
}

func (n Npc) Style() tcell.Style {
	return friendlyNPCStyle
}

func (n Npc) Collides(with Position) bool {
	return collides(n.Position(), with)
}

var friendlyNPCStyle = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorGreen).Bold(true)
