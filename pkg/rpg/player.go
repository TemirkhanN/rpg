package rpg

type Player struct {
	name        string
	whereabouts Location
}

func NewPlayer(playerName string, location Location) Player {
	return Player{
		name:        playerName,
		whereabouts: location,
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

func (p Player) StartConversation(npc NPC) string {
	return "Hello " + p.Name() + " I am " + npc.Name()
}
