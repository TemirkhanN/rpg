package game

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/pkg/rpg"
)

func New() {
	newNPC := rpg.NewNPC("Ded")
	newPlayer := rpg.NewPlayer("Nova")

	locations := []rpg.Location{
		rpg.NewLocation("Gludio"),
		rpg.NewLocation("Dion"),
		rpg.NewLocation("Goddard"),
	}

	fmt.Println(newPlayer.StartConversation(newNPC))

	screen := createScreen()
	screen.EnableMouse()

	for {
		draw(screen, &newPlayer, locations)
	}
}

func draw(screen tcell.Screen, player *rpg.Player, locations []rpg.Location) {
	screen.Clear()

	width, _ := screen.Size()
	screenCenter := width / 2
	newText("RPG maker", screenCenter-10, 1).draw(screen, textStyle)

	newBox(1, 1, 25, 10, "Status").draw(screen, boxStyle)

	newText("Name: "+player.Name(), 2, 3).draw(screen, infoTextStyle)
	newText("Currently at: "+player.Whereabouts().Name(), 2, 4).draw(screen, infoTextStyle)

	teleport := newBox(1, 12, 25, 20, "Move to")
	teleport.draw(screen, boxStyle)

	for i, location := range locations {
		newText(
			fmt.Sprintf("%d.%s", i+1, location.Name()),
			teleport.leftTop.x+1,
			teleport.leftTop.y+i+1,
		).draw(screen, infoTextStyle)
	}

	totalLocationsAmount := len(locations)

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

			if (mouseX >= teleport.leftTop.x+1 && mouseX <= teleport.rightBottom.x-1) &&
				(mouseY >= teleport.leftTop.y+1 && mouseY <= teleport.leftTop.y+totalLocationsAmount) {
				locationKey := mouseY - teleport.leftTop.y - 1
				if player.Whereabouts().Name() != locations[locationKey].Name() {
					player.MoveToLocation(locations[locationKey])
				}
			}
		}
	}

	screen.Show()
}
