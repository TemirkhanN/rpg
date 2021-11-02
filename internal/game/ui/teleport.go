package ui

import (
	"fmt"

	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type Teleport struct {
	locations []rpg.Location
	panel     Box
}

func NewTeleport(locations []rpg.Location, panel Box) Teleport {
	return Teleport{
		locations: locations,
		panel:     panel,
	}
}

func (t Teleport) Click(player *rpg.Player, pos Position) {
	leftTop := t.panel.leftTop
	for i, location := range t.locations {
		optionPos := Position{leftTop.X + 1, leftTop.Y + i + 1}
		// if option is not on the same line or Text is not within clicked Position.
		if pos.Y != optionPos.Y ||
			pos.X < optionPos.X ||
			pos.X > optionPos.X+CalculateTextWidth(location.Name()) {
			continue
		}

		if player.Whereabouts().Name() != location.Name() {
			player.MoveToLocation(location)
		}

		break
	}
}

func (t Teleport) Draw(on tcell.Screen, panelStyle tcell.Style, optionsStyle tcell.Style) {
	t.panel.Draw(on, panelStyle)

	leftTop := t.panel.leftTop
	for i, location := range t.locations {
		optionPos := Position{leftTop.X + 1, leftTop.Y + i + 1}
		NewText(fmt.Sprintf("* %s", location.Name()), optionPos.X, optionPos.Y).Draw(on, optionsStyle)
	}
}
