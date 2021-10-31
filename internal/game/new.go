package game

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type Game struct {
	player    *rpg.Player
	locations []rpg.Location

	title    *text
	status   *playerPanel
	teleport *teleportPanel
}

type playerPanel struct {
	panel    *box
	name     *text
	location *text
}

type teleportPanel struct {
	panel     *box
	locations []text
}

func New(playerName string) *Game {
	newGame := new(Game)
	newPlayer := rpg.NewPlayer(playerName)
	newGame.player = &newPlayer

	newGame.locations = []rpg.Location{
		rpg.NewLocation("Gludio"),
		rpg.NewLocation("Dion"),
		rpg.NewLocation("Goddard"),
	}

	newNPC := rpg.NewNPC("Ded")
	fmt.Println(newPlayer.StartConversation(newNPC))

	return newGame
}

func (g *Game) Start() {
	screen := createScreen()
	screen.EnableMouse()
	width, _ := screen.Size()
	screenCenter := width / 2

	g.title = newText("RPG maker", screenCenter-10, 1)
	g.title = newText("RPG maker", screenCenter-10, 1)
	g.status = &playerPanel{
		panel:    newBox(1, 1, 25, 10, "Status"),
		name:     newText("Name: "+g.player.Name(), 2, 3),
		location: newText("Currently at: "+g.player.Whereabouts().Name(), 2, 4),
	}

	g.teleport = &teleportPanel{
		panel:     newBox(1, 12, 25, 20, "Move to"),
		locations: make([]text, len(g.locations)),
	}

	for i, location := range g.locations {
		g.teleport.locations[i] = *newText(
			fmt.Sprintf("%d.%s", i+1, location.Name()),
			g.teleport.panel.leftTop.x+1,
			g.teleport.panel.leftTop.y+i+1,
		)
	}

	g.drawGraphics(screen)
	for {
		g.handleEvents(screen)
	}
}

func (g Game) drawGraphics(screen tcell.Screen) {
	screen.Clear()

	g.title.draw(screen, infoTextStyle)

	// Player status panel.
	g.status.panel.draw(screen, boxStyle)
	g.status.name.draw(screen, infoTextStyle)
	g.status.location.draw(screen, infoTextStyle)

	// Teleport panel.
	g.teleport.panel.draw(screen, boxStyle)
	for _, location := range g.teleport.locations {
		location.draw(screen, infoTextStyle)
	}

	screen.Show()
}

func (g Game) handleEvents(screen tcell.Screen) {
	switch ev := screen.PollEvent().(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape {
			screen.Fini()
			os.Exit(0)
		}
	case *tcell.EventMouse:
		// left click
		if ev.Buttons() == tcell.Button1 {
			mouseX, mouseY := ev.Position()
			teleport := g.teleport.panel

			if (mouseX >= teleport.leftTop.x+1 && mouseX <= teleport.rightBottom.x-1) &&
				(mouseY >= teleport.leftTop.y+1 && mouseY <= teleport.leftTop.y+len(g.locations)) {
				locationKey := mouseY - teleport.leftTop.y - 1
				if g.player.Whereabouts().Name() != g.locations[locationKey].Name() {
					g.player.MoveToLocation(g.locations[locationKey])
					g.status.location.text = fmt.Sprintf("Currently at: %s", g.player.Whereabouts().Name())

					g.drawGraphics(screen)
				}
			}
		}
	}
}
