package scenes

import (
	"os"
	"strconv"

	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/internal/game/data"
	"github.com/TemirkhanN/rpg/internal/game/resources"
	"github.com/TemirkhanN/rpg/internal/game/ui"
	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type Scene interface {
	Draw(on tcell.Screen)
	HandleEvents(within tcell.Screen)
}

type MainScene struct {
	player *data.Player

	res resources.Resources
}

func NewMainScene(resources resources.Resources, player *rpg.Player) MainScene {
	mainScene := MainScene{
		player: data.NewPlayer(player),
		res:    resources,
	}

	location, err := resources.LoadLocation(1)
	if err != nil {
		panic(err)
	}

	mainScene.player.TeleportToLocation(location, data.NewPos(35, 10))

	return mainScene
}

func (ms MainScene) Draw(on tcell.Screen) {
	on.Clear()

	location := ms.player.Whereabouts()
	ui.DrawLocation(location, on)
	ui.DrawPlayer(*ms.player, on)
	ui.DrawStatusPanel(*ms.player, on)

	on.Show()
}

func (ms MainScene) HandleEvents(within tcell.Screen) {
	switch ev := within.PollEvent().(type) {
	case *tcell.EventKey:
		ms.handleControlsInput(within, *ev)
	case *tcell.EventMouse:
		/*
			// left click
			if ev.Buttons() == tcell.Button1 {
			}
		*/
	}
}

func (ms MainScene) handleControlsInput(within tcell.Screen, input tcell.EventKey) {
	key := input.Key()
	// Close game.
	if key == tcell.KeyEscape {
		within.Fini()
		os.Exit(0)
	}

	choice, _ := strconv.Atoi(string(input.Rune()))
	if choice > 0 {
		ms.player.ChooseDialogueReplyOption(choice)

		return
	}

	ms.handleMovement(key)
}

func (ms MainScene) handleMovement(key tcell.Key) {
	newPosition := ms.player.Position()
	if key == tcell.KeyUp {
		newPosition = data.NewPos(newPosition.X(), newPosition.Y()-1)
	}

	if key == tcell.KeyDown {
		newPosition = data.NewPos(newPosition.X(), newPosition.Y()+1)
	}

	if key == tcell.KeyLeft {
		newPosition = data.NewPos(newPosition.X()-1, newPosition.Y())
	}

	if key == tcell.KeyRight {
		newPosition = data.NewPos(newPosition.X()+1, newPosition.Y())
	}

	if newPosition != ms.player.Position() {
		ms.player.MoveTo(newPosition)
	}
}
