package game

import (
	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/internal/game/resources"
	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type window interface {
	draw()
	handleEvents()
}

type Game struct {
	player    *player
	resources resources.Resources

	screen       tcell.Screen
	activeWindow window
}

func New(playerName string) *Game {
	newGame := new(Game)
	newGame.resources = resources.LoadResources()
	newbieTown, _ := newGame.resources.GetLocation("Talking Island")
	newPlayer := rpg.NewPlayer(playerName, newbieTown)
	newGame.player = &player{
		asci: asci{
			symbol: '🐶',
			style:  playerStyle,
		},
		currentDialogue:     rpg.NoDialogue,
		currentDialogueWith: rpg.NoNpc,
		player:              &newPlayer,
		pos:                 position{x: 30, y: 15},
	}

	newGame.screen = createScreen()
	newGame.screen.EnableMouse()

	newGame.activeWindow = newPrimaryWindow(*newGame)

	return newGame
}

func (g *Game) Start() {
	g.activeWindow.draw()
	for {
		g.activeWindow.handleEvents()
	}
}
