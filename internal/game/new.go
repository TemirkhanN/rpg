package game

import (
	"os"
	"unicode/utf8"

	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type Game struct {
	player   *player
	npcs     []*npc
	teleport teleport

	screen tcell.Screen
}

func New(playerName string) *Game {
	newGame := new(Game)
	newbieTown := rpg.NewLocation("Neu Beetown")
	newPlayer := rpg.NewPlayer(playerName, newbieTown)
	newGame.player = &player{
		asci: asci{
			symbol: '🐶',
			style:  playerStyle,
		},
		currentDialogue: "",
		player:          &newPlayer,
		pos:             position{x: 30, y: 15},
	}

	newGame.teleport = newTeleport(
		[]rpg.Location{
			newbieTown,
			rpg.NewLocation("Gludio"),
			rpg.NewLocation("Dion"),
			rpg.NewLocation("Goddard"),
		},
		*newBox(1, 12, 25, 20, "Move to"),
	)

	dedNPC := rpg.NewNPC("Ded", map[string]rpg.Dialogue{
		"defaultDialogue": rpg.NewDialogue("What's your profession?",
			[]string{
				"I'm a warrior",
				"I'm a mage",
			}),
	})

	newbieHelper := rpg.NewNPC("Newbie Helper", map[string]rpg.Dialogue{
		"defaultDialogue": rpg.NewDialogue("I can help you",
			[]string{
				"Do you need help?",
			}),
	})
	newGame.npcs = []*npc{
		{
			asci: asci{style: friendlyNPCStyle, symbol: '👱'},
			npc:  newbieHelper,
			pos:  position{x: 40, y: 8},
		},
		{
			asci: asci{style: friendlyNPCStyle, symbol: '🧔'},
			npc:  dedNPC,
			pos:  position{x: 55, y: 13},
		},
	}

	newGame.screen = createScreen()
	newGame.screen.EnableMouse()

	return newGame
}

func (g *Game) Start() {
	g.drawGraphics()
	for {
		g.handleEvents()
	}
}

func (g Game) drawGraphics() {
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
	g.teleport.draw(g.screen, boxStyle, infoTextStyle)

	// Draw npcs.
	for _, gameNpc := range g.npcs {
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

func (g Game) handleEvents() {
	switch ev := g.screen.PollEvent().(type) {
	case *tcell.EventKey:
		g.handleControlsInput(ev.Key())
	case *tcell.EventMouse:
		// left click
		if ev.Buttons() == tcell.Button1 {
			mouseX, mouseY := ev.Position()

			// Handle click on teleport panel.
			g.teleport.click(g.player.player, position{x: mouseX, y: mouseY})

			g.drawGraphics()
		}
	}
}

func (g Game) handleControlsInput(input tcell.Key) {
	if input == tcell.KeyEscape {
		g.screen.Fini()
		os.Exit(0)
	}

	if input == tcell.KeyUp {
		newPosition := g.player.pos
		newPosition.y--

		g.playerMoveTo(newPosition)
	}

	if input == tcell.KeyDown {
		newPosition := g.player.pos
		newPosition.y++

		g.playerMoveTo(newPosition)
	}

	if input == tcell.KeyLeft {
		newPosition := g.player.pos
		newPosition.x--

		g.playerMoveTo(newPosition)
	}

	if input == tcell.KeyRight {
		newPosition := g.player.pos
		newPosition.x++

		g.playerMoveTo(newPosition)
	}
}

func (g Game) playerMoveTo(to position) {
	movementAllowed := true
	isInDialogue := false

	for _, npc := range g.npcs {
		if npc.collides(to) {
			dialogue := g.player.player.StartConversation(npc.npc)
			g.player.currentDialogue = dialogue.Text()
			movementAllowed = false
			isInDialogue = true
			g.drawGraphics()

			break
		}
	}

	if !isInDialogue && g.player.currentDialogue != "" {
		g.player.currentDialogue = ""
	}

	if movementAllowed {
		g.player.pos = to
		g.drawGraphics()
	}
}
