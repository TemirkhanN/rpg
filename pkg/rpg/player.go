package rpg

type Player struct {
	name        string
	whereabouts Location
}

func NewPlayer(playerName string) Player {
	return Player{
		name:        playerName,
		whereabouts: NewLocation("Unknown location"),
	}
}

func (p Player) Name() string {
	return p.name
}

func (p Player) Whereabouts() Location {
	return p.whereabouts
}

func (p *Player) MoveToLocation(location Location) {
	p.whereabouts = location
}
