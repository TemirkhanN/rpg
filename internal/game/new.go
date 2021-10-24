package game

import (
	"fmt"

	"github.com/TemirkhanN/rpg/pkg/rpg"
)

func New() {
	newPlayer := rpg.NewPlayer("Nova")
	fmt.Println(newPlayer.Name())
	newLocation := rpg.NewLocation("Gludio")
	fmt.Println(newLocation.Name())
}
