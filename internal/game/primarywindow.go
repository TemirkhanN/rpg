package game

import (
	"os"
	"unicode/utf8"

	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type primaryWindow struct {
	game     Game
	teleport teleport
	npcs     []*npc
}

func newPrimaryWindow(game Game) primaryWindow {
	return primaryWindow{
		game: game,
		teleport: newTeleport(
			[]rpg.Location{
				game.player.player.Whereabouts(),
				rpg.NewLocation("Gludio"),
				rpg.NewLocation("Dion"),
				rpg.NewLocation("Goddard"),
			},
			*newBox(1, 12, 25, 20, "Move to"),
		),
		npcs: []*npc{
			{
				asci: asci{style: friendlyNPCStyle, symbol: '👱'},
				npc: rpg.NewNPC("Ded", map[string]rpg.Dialogue{
					"defaultDialogue": rpg.NewDialogue("What's your profession?",
						[]string{
							"I'm a warrior",
							"I'm a mage",
						}),
				}),
				pos: position{x: 40, y: 8},
			},
			{
				asci: asci{style: friendlyNPCStyle, symbol: '🧔'},
				npc: rpg.NewNPC("Newbie Helper", map[string]rpg.Dialogue{
					"defaultDialogue": rpg.NewDialogue("I can help you",
						[]string{
							"Do you need help?",
						}),
				}),
				pos: position{x: 55, y: 13},
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
	if g.player.currentDialogue != "" {
		newBox(80, 12, 110, 20, "Dialogue").draw(g.screen, boxStyle)
		verticalOffset := 13

		charactersPerLine := 110 - 80 - 4
		dialogueLength := utf8.RuneCountInString(g.player.currentDialogue)
		if dialogueLength <= charactersPerLine {
			newText(g.player.currentDialogue, 82, verticalOffset).draw(g.screen, textStyle)
		} else {
			var linedString string
			for i, r := range g.player.currentDialogue {
				linedString += string(r)

				if (i+1)%charactersPerLine == 0 || i+1 == dialogueLength {
					newText(linedString, 82, verticalOffset).draw(g.screen, textStyle)
					verticalOffset++
					linedString = ""
				}
			}
		}
	}

	g.screen.Show()
}

func (pw primaryWindow) handleEvents() {
	g := pw.game
	switch ev := g.screen.PollEvent().(type) {
	case *tcell.EventKey:
		pw.handleControlsInput(ev.Key())
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

func (pw primaryWindow) handleControlsInput(input tcell.Key) {
	g := pw.game

	if input == tcell.KeyEscape {
		g.screen.Fini()
		os.Exit(0)
	}

	if input == tcell.KeyUp {
		newPosition := g.player.pos
		newPosition.y--

		pw.playerMoveTo(newPosition)
	}

	if input == tcell.KeyDown {
		newPosition := g.player.pos
		newPosition.y++

		pw.playerMoveTo(newPosition)
	}

	if input == tcell.KeyLeft {
		newPosition := g.player.pos
		newPosition.x--

		pw.playerMoveTo(newPosition)
	}

	if input == tcell.KeyRight {
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
			dialogue := g.player.player.StartConversation(npc.npc)
			g.player.currentDialogue = dialogue.Text()
			movementAllowed = false
			isInDialogue = true
			pw.draw()

			break
		}
	}

	if !isInDialogue && g.player.currentDialogue != "" {
		g.player.currentDialogue = ""
	}

	if movementAllowed {
		g.player.pos = to
		pw.draw()
	}
}
