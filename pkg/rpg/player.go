package rpg

type Player struct {
	name           string
	whereabouts    Location
	activeDialogue Dialogue
}

func NewPlayer(playerName string) Player {
	return Player{
		name: playerName,
		whereabouts: Location{
			id:   0,
			name: "",
		},
		activeDialogue: NoDialogue,
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

func (p *Player) StartConversation(npc NPC) Dialogue {
	p.activeDialogue = npc.StartConversation(*p)

	return p.ActiveDialogue()
}

func (p *Player) EndConversation() {
	p.activeDialogue = NoDialogue
}

func (p Player) ActiveDialogue() Dialogue {
	return p.activeDialogue
}

func (p *Player) Reply(to NPC, with Phrase) Dialogue {
	p.activeDialogue = to.Reply(with)

	return p.ActiveDialogue()
}
