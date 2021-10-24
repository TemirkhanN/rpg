package rpg

type Player struct {
	name string
}

func NewPlayer(playerName string) Player {
	return Player{
		name: playerName,
	}
}

func (p Player) Name() string {
	return p.name
}
