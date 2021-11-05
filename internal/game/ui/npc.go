package ui

import (
	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type NPC struct {
	asci asci
	npc  *rpg.NPC
	pos  Position
}

func NewNPC(npc *rpg.NPC, pos Position, icon rune, iconStyle tcell.Style) NPC {
	return NPC{
		asci: asci{style: iconStyle, symbol: icon},
		npc:  npc,
		pos:  pos,
	}
}

func (n NPC) Draw(on tcell.Screen) {
	NewText(string(n.asci.symbol), n.pos.X, n.pos.Y).Draw(on, n.asci.style)
}

func (n NPC) Collides(with Position) bool {
	if with.Y != n.pos.Y {
		return false
	}

	// npc takes 3 slots on horizontal line
	if with.X < n.pos.X-1 || with.X > n.pos.X+1 {
		return false
	}

	return true
}

var (
	FriendlyNPCStyle = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorGreen).Bold(true)
	NeutralNPCStyle  = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen).Bold(true)
)
