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

func New(playerName string) Game {
	newPlayer := rpg.NewPlayer(playerName)
	screen := ui.CreateScreen()
	screen.EnableMouse()

	newGame := Game{
		player:      &newPlayer,
		resources:   resources.LoadResources(),
		screen:      screen,
		activeScene: nil,
	}

	newGame.activeScene = scenes.NewMainScene(newGame.resources, newGame.player)

	return newGame
}

func (g *Game) Start() {
	for {
		g.activeScene.Draw(g.screen)
		g.activeScene.HandleEvents(g.screen)
	}
}
