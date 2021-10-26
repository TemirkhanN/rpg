package game

import (
	"fmt"

	"github.com/TemirkhanN/rpg/pkg/rpg"
)

func New() {
	newNPC := rpg.NewNPC("Ded")
	fmt.Println("Created NPC: " + newNPC.Name())

	newPlayer := rpg.NewPlayer("Nova")
	fmt.Println("Created player: " + newPlayer.Name())

	Gludio := rpg.NewLocation("Gludio")
	fmt.Println("Created location: " + Gludio.Name())

	fmt.Println(newPlayer.Name() + " is currently in " + newPlayer.Whereabouts().Name())
	newPlayer.MoveToLocation(Gludio)
	fmt.Println(newPlayer.Name() + " teleported to " + newPlayer.Whereabouts().Name())

	fmt.Println(newPlayer.StartConversation(newNPC))
}
