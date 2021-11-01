package game

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type primaryWindow struct {
	game     Game
	teleport teleport
	npcs     []*npc
}

func newPrimaryWindow(game Game) primaryWindow {
	newbieHelper, _ := game.resources.GetNPC("Newbie Helper")
	guard, _ := game.resources.GetNPC("Guard")

	return primaryWindow{
		game: game,
		teleport: newTeleport(
			game.resources.Locations(),
			*newBox(1, 12, 25, 20, "Move to"),
		),
		npcs: []*npc{
			{
				asci: asci{style: friendlyNPCStyle, symbol: '👱'},
				npc:  newbieHelper,
				pos:  position{x: 40, y: 8},
			},
			{
				asci: asci{style: friendlyNPCStyle, symbol: '🧔'},
				npc:  guard,
				pos:  position{x: 55, y: 13},
			},
		},
	}
}

func (pw primaryWindow) draw() {
	g := pw.game

	g.screen.Clear()

	width, _ := g.screen.Size()
	screenCenter := width / 2

	// Game title.
	newText("RPG maker", screenCenter-10, 1).draw(g.screen, infoTextStyle)

	// Status panel.
	newBox(1, 1, 25, 10, "Status").draw(g.screen, boxStyle)
	newText("Name: "+g.player.player.Name(), 2, 3).draw(g.screen, infoTextStyle)
	newText("Location: "+g.player.player.Whereabouts().Name(), 2, 4).draw(g.screen, infoTextStyle)

	// Teleport panel.
	pw.teleport.draw(g.screen, boxStyle, infoTextStyle)

	// Draw npcs.
	for _, gameNpc := range pw.npcs {
		newText(string(gameNpc.asci.symbol), gameNpc.pos.x, gameNpc.pos.y).draw(g.screen, gameNpc.asci.style)
	}

	// Draw player.
	newText(string(g.player.asci.symbol), g.player.pos.x, g.player.pos.y).draw(g.screen, g.player.asci.style)

	// Draw player dialogue panel.
	if !g.player.currentDialogue.Empty() {
		newBox(80, 12, 110, 20, "Dialogue").draw(g.screen, boxStyle)
		at := position{x: 80 + 2, y: 12 + 1}
		charactersPerLine := 110 - 80 - 4

		textEndsAt := drawTextWithAutoLinebreaks(g.screen, g.player.currentDialogue.Text(), at, charactersPerLine)

		for i, reply := range g.player.currentDialogue.Choices() {
			textEndsAt = drawTextWithAutoLinebreaks(
				g.screen,
				fmt.Sprintf("%d.%s", i+1, reply),
				textEndsAt,
				charactersPerLine,
			)
		}
	}

	g.screen.Show()
}

func (pw primaryWindow) handleEvents() {
	g := pw.game
	switch ev := g.screen.PollEvent().(type) {
	case *tcell.EventKey:
		pw.handleControlsInput(*ev)
	case *tcell.EventMouse:
		// left click
		if ev.Buttons() == tcell.Button1 {
			mouseX, mouseY := ev.Position()

			// Handle click on teleport panel.
			pw.teleport.click(g.player.player, position{x: mouseX, y: mouseY})

			pw.draw()
		}
	}
}

func (pw primaryWindow) handleControlsInput(input tcell.EventKey) {
	key := input.Key()
	// Close game.
	if key == tcell.KeyEscape {
		pw.game.screen.Fini()
		os.Exit(0)
	}

	pw.handleDialogueReply(input)
	pw.handleMovement(key)
}

func (pw primaryWindow) handleDialogueReply(input tcell.EventKey) {
	g := pw.game

	currentDialogue := g.player.currentDialogue
	if currentDialogue.Empty() {
		return
	}

	availableReplies := currentDialogue.Choices()
	if len(availableReplies) == 0 {
		return
	}

	choice, err := strconv.Atoi(string(input.Rune()))
	if err != nil || choice < 0 || len(availableReplies) < choice {
		return
	}

	dialogue := pw.game.player.player.Reply(pw.game.player.currentDialogueWith, availableReplies[choice-1])
	if dialogue.Text() != g.player.currentDialogue.Text() {
		g.player.currentDialogue = dialogue
		pw.draw()
	}
}

func (pw primaryWindow) handleMovement(key tcell.Key) {
	g := pw.game
	if key == tcell.KeyUp {
		newPosition := g.player.pos
		newPosition.y--

		pw.playerMoveTo(newPosition)
	}

	if key == tcell.KeyDown {
		newPosition := g.player.pos
		newPosition.y++

		pw.playerMoveTo(newPosition)
	}

	if key == tcell.KeyLeft {
		newPosition := g.player.pos
		newPosition.x--

		pw.playerMoveTo(newPosition)
	}

	if key == tcell.KeyRight {
		newPosition := g.player.pos
		newPosition.x++

		pw.playerMoveTo(newPosition)
	}
}

func (pw primaryWindow) playerMoveTo(to position) {
	g := pw.game

	movementAllowed := true
	isInDialogue := false

	for _, npc := range pw.npcs {
		if npc.collides(to) {
			movementAllowed = false
			isInDialogue = true
			dialogue := g.player.player.StartConversation(npc.npc)
			if dialogue.Text() == g.player.currentDialogue.Text() {
				break
			}
			g.player.currentDialogue = dialogue
			g.player.currentDialogueWith = npc.npc
			pw.draw()

			break
		}
	}

	if !isInDialogue && !g.player.currentDialogue.Empty() {
		g.player.currentDialogue = rpg.NoDialogue
		g.player.currentDialogueWith = rpg.NoNpc
	}

	if movementAllowed {
		g.player.pos = to
		pw.draw()
	}
}
