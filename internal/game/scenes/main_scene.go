package scenes

import (
	"os"

	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/internal/game/resources"
	"github.com/TemirkhanN/rpg/internal/game/ui"
	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type Scene interface {
	Draw(on tcell.Screen)
	HandleEvents(within tcell.Screen)
}

type MainScene struct {
	teleport ui.Teleport
	npcs     []ui.NPC
	player   *ui.Player
}

func NewMainScene(resources resources.Resources, player *rpg.Player) MainScene {
	newbieHelper, _ := resources.GetNPC("Newbie Helper")
	guard, _ := resources.GetNPC("Guard")

	playerUI := ui.NewPlayer(player, ui.Position{X: 30, Y: 15}, '🐶', ui.PlayerIconStyle)

	return MainScene{
		teleport: ui.NewTeleport(
			resources.Locations(),
			*ui.NewBox(1, 12, 25, 20, "Move to"),
		),
		player: &playerUI,
		npcs: []ui.NPC{
			ui.NewNPC(&newbieHelper, ui.Position{X: 40, Y: 8}, '👱', ui.FriendlyNPCStyle),
			ui.NewNPC(&guard, ui.Position{X: 55, Y: 13}, '🧔', ui.FriendlyNPCStyle),
		},
	}
}

func (ms MainScene) Draw(on tcell.Screen) {
	on.Clear()

	width, _ := on.Size()
	screenCenter := width / 2

	// Game title.
	ui.NewText("RPG maker", screenCenter-10, 1).Draw(on, ui.InfoTextStyle)

	// Status panel.
	ui.NewBox(1, 1, 25, 10, "Status").Draw(on)
	ui.NewText("Name: "+ms.player.Player().Name(), 2, 3).Draw(on, ui.InfoTextStyle)
	ui.NewText("Location: "+ms.player.Player().Whereabouts().Name(), 2, 4).Draw(on, ui.InfoTextStyle)

	// Teleport panel.
	ms.teleport.Draw(on, ui.BoxStyle, ui.InfoTextStyle)

	// Draw npcs.
	for _, gameNpc := range ms.npcs {
		gameNpc.Draw(on)
	}

	// Draw player.
	ms.player.Draw(on)

	// Draw player dialogue panel.
	ms.player.DrawDialogue(on, ui.Position{X: 80, Y: 1})

	on.Show()
}

func (ms MainScene) HandleEvents(within tcell.Screen) {
	switch ev := within.PollEvent().(type) {
	case *tcell.EventKey:
		ms.handleControlsInput(within, *ev)
	case *tcell.EventMouse:
		// left click
		if ev.Buttons() == tcell.Button1 {
			mouseX, mouseY := ev.Position()

			// Handle click on teleport panel.
			ms.teleport.Click(ms.player.Player(), ui.Position{X: mouseX, Y: mouseY})

			ms.Draw(within)
		}
	}
}

func (ms MainScene) handleControlsInput(within tcell.Screen, input tcell.EventKey) {
	key := input.Key()
	// Close game.
	if key == tcell.KeyEscape {
		within.Fini()
		os.Exit(0)
	}

	ms.player.ChooseDialogueReplyOption(input)
	ms.handleMovement(key)

	ms.Draw(within)
}

func (ms MainScene) handleMovement(key tcell.Key) {
	if key == tcell.KeyUp {
		newPosition := ms.player.Position()
		newPosition.Y--

		ms.playerMoveTo(newPosition)
	}

	if key == tcell.KeyDown {
		newPosition := ms.player.Position()
		newPosition.Y++

		ms.playerMoveTo(newPosition)
	}

	if key == tcell.KeyLeft {
		newPosition := ms.player.Position()
		newPosition.X--

		ms.playerMoveTo(newPosition)
	}

	if key == tcell.KeyRight {
		newPosition := ms.player.Position()
		newPosition.X++

		ms.playerMoveTo(newPosition)
	}
}

func (ms MainScene) playerMoveTo(to ui.Position) {
	movementAllowed := true
	isInDialogue := false

	for _, npc := range ms.npcs {
		if npc.Collides(to) {
			movementAllowed = false
			isInDialogue = true
			ms.player.StartDialogue(npc)

			break
		}
	}

	if !isInDialogue {
		ms.player.EndDialogue()
	}

	if movementAllowed {
		ms.player.MoveTo(to)
	}
}
