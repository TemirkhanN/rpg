package ui

import (
	"github.com/gdamore/tcell"
	"github.com/gookit/event"

	"github.com/TemirkhanN/rpg/internal/game/ui/icon"
	"github.com/TemirkhanN/rpg/internal/game/ui/interactive"
	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type Door struct {
	SolidObject
	interactive.Element
}

func NewDoor(to rpg.Location, pos Position) Door {
	doorColor := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorSandyBrown).Bold(true)

	return Door{
		SolidObject: NewSolidObject(icon.Door, doorColor, pos),
		Element: interactive.NewElement(func() {
			_ = event.TriggerEvent(event.NewBasic("OpenedDoor", event.M{"to": to}))
		}),
	}
}
