package game

import (
	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/internal/game/resources"
	"github.com/TemirkhanN/rpg/internal/game/scenes"
	"github.com/TemirkhanN/rpg/internal/game/ui"
	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type Game struct {
	player    *rpg.Player
	resources resources.Resources

	screen      tcell.Screen
	activeScene scenes.Scene
}

func (g Game) Screen() tcell.Screen {
	return g.screen
}

func New(playerName string) *Game {
	newGame := new(Game)
	newGame.resources = resources.LoadResources()
	newbieTown, _ := newGame.resources.GetLocation("Talking Island")
	newPlayer := rpg.NewPlayer(playerName, newbieTown)
	newGame.player = &newPlayer

	newGame.screen = ui.CreateScreen()
	newGame.screen.EnableMouse()

	newGame.activeScene = scenes.NewMainScene(newGame.resources, newGame.player)

	return newGame
}

func (g *Game) Start() {
	g.activeScene.Draw(g.screen)
	for {
		g.activeScene.HandleEvents(g.screen)
	}
}
