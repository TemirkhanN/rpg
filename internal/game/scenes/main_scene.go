package scenes

import (
	"os"

	"github.com/gdamore/tcell"
	"github.com/gookit/event"

	"github.com/TemirkhanN/rpg/internal/game/resources"
	"github.com/TemirkhanN/rpg/internal/game/ui"
	"github.com/TemirkhanN/rpg/internal/game/ui/icon"
	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type Scene interface {
	Draw(on tcell.Screen)
	HandleEvents(within tcell.Screen)
}

type MainScene struct {
	npcs       []ui.NPC
	player     *ui.Player
	area       ui.Area
	doors      []ui.Door
	solidItems []ui.SolidObject

	res resources.Resources
}

func NewMainScene(resources resources.Resources, player *rpg.Player) *MainScene {
	playerUI := ui.NewPlayer(player, ui.Position{X: 0, Y: 0}, icon.DogFace, ui.PlayerIconStyle)

	mainScene := &MainScene{
		npcs:       nil,
		player:     &playerUI,
		area:       ui.Area{},
		doors:      nil,
		solidItems: nil,
		res:        resources,
	}

	event.On("OpenedDoor", event.ListenerFunc(func(e event.Event) error {
		location := e.Get("to").(rpg.Location)

		player.MoveToLocation(location)
		mainScene.loadLocation(player.Whereabouts())

		return nil
	}))

	mainScene.loadTalkingIsland()

	mainScene.player.Teleport(ui.Position{X: 35, Y: 10})

	return mainScene
}

// todo create some dynamic preset for different locations.
func (ms *MainScene) loadLocation(location rpg.Location) {
	switch location.Name() {
	case "Talking Island":
		ms.loadTalkingIsland()
	case "Forest":
		ms.loadForest()
	}
}

func (ms *MainScene) clearScene() {
	ms.npcs = nil
	ms.doors = nil
	ms.area = ui.Area{}
	ms.solidItems = nil
}

func (ms *MainScene) loadTalkingIsland() {
	ms.clearScene()

	newbieHelper, _ := ms.res.GetNPC("Newbie Helper")
	guard, _ := ms.res.GetNPC("Guard")
	cow, _ := ms.res.GetNPC("Cow")

	ms.npcs = []ui.NPC{
		ui.NewNPC(&newbieHelper, ui.Position{X: 40, Y: 8}, icon.WomanFace, ui.FriendlyNPCStyle),
		ui.NewNPC(&guard, ui.Position{X: 55, Y: 13}, icon.ManFace, ui.FriendlyNPCStyle),
		ui.NewNPC(&cow, ui.Position{X: 50, Y: 17}, icon.Cow, ui.NeutralNPCStyle),
	}

	ms.area = ui.NewArea(26, 1, 79, 20, ui.BoxStyle)

	forest, _ := ms.res.GetLocation("Forest")
	ms.doors = []ui.Door{
		ui.NewDoor(forest, ui.Position{X: 41, Y: 1}),
	}

	ms.player.Teleport(ui.Position{X: 41, Y: 2})
}

func (ms *MainScene) loadForest() {
	ms.clearScene()

	guard, _ := ms.res.GetNPC("Guard")
	cow1, _ := ms.res.GetNPC("Cow")
	cow2, _ := ms.res.GetNPC("Cow")
	cow3, _ := ms.res.GetNPC("Cow")
	cow4, _ := ms.res.GetNPC("Cow")
	cow5, _ := ms.res.GetNPC("Cow")

	ms.npcs = []ui.NPC{
		ui.NewNPC(&guard, ui.Position{X: 35, Y: 19}, icon.ManFace, ui.FriendlyNPCStyle),
		ui.NewNPC(&cow1, ui.Position{X: 30, Y: 12}, icon.Cow, ui.NeutralNPCStyle),
		ui.NewNPC(&cow2, ui.Position{X: 36, Y: 10}, icon.Cow, ui.NeutralNPCStyle),
		ui.NewNPC(&cow3, ui.Position{X: 40, Y: 5}, icon.Cow, ui.NeutralNPCStyle),
		ui.NewNPC(&cow4, ui.Position{X: 32, Y: 5}, icon.Cow, ui.NeutralNPCStyle),
		ui.NewNPC(&cow5, ui.Position{X: 37, Y: 13}, icon.Cow, ui.NeutralNPCStyle),
	}

	ms.area = ui.NewArea(26, 1, 79, 20, ui.BoxStyle)

	talkingIsland, _ := ms.res.GetLocation("Talking Island")
	ms.doors = []ui.Door{
		ui.NewDoor(talkingIsland, ui.Position{X: 41, Y: 20}),
	}

	treeColor := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorForestGreen).Bold(true)
	for x := 28; x < 78; x += 3 {
		for y := 2; y < 5; y++ {
			ms.solidItems = append(
				ms.solidItems,
				ui.NewSolidObject(icon.EvergreenTree, treeColor, ui.Position{X: x, Y: y}),
			)
		}
	}

	ms.player.Teleport(ui.Position{X: 41, Y: 19})
}

func (ms MainScene) Draw(on tcell.Screen) {
	on.Clear()

	// Status panel.
	ui.NewBox(1, 1, 25, 10, "Status").Draw(on)
	ui.NewText("Name: "+ms.player.Player().Name(), 2, 3).Draw(on, ui.InfoTextStyle)
	ui.NewText("Location: "+ms.player.Player().Whereabouts().Name(), 2, 4).Draw(on, ui.InfoTextStyle)

	// Draw npcs.
	for _, gameNpc := range ms.npcs {
		gameNpc.Draw(on)
	}

	// Draw player.
	ms.player.Draw(on)

	// Draw player dialogue panel.
	ms.player.DrawDialogue(on, ui.Position{X: 80, Y: 1})

	ms.area.Draw(on)

	for _, door := range ms.doors {
		door.Draw(on)
	}

	for _, item := range ms.solidItems {
		item.Draw(on)
	}

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
	for _, solidItem := range ms.solidItems {
		if solidItem.Collides(to) {
			return
		}
	}

	for _, door := range ms.doors {
		if door.Collides(to) {
			ms.player.Interact(door.Element)

			return
		}
	}

	movementAllowed := true
	isInDialogue := false

	if !to.Inside(ms.area) {
		return
	}

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
