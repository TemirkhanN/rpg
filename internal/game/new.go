package game

import (
	"fmt"

	"github.com/TemirkhanN/rpg/pkg/rpg"
)

func New() {
	newPlayer := rpg.NewPlayer("Nova")
	newPlayer2 := rpg.NewPlayer("Bozman")
	fmt.Println(newPlayer.Name())
	fmt.Println(newPlayer2.Name())
}
