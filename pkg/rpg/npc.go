package rpg

type NPC struct {
	name      string
	dialogues map[string]Dialogue
}

func NewNPC(name string, dialogues map[string]Dialogue) NPC {
	return NPC{
		name:      name,
		dialogues: dialogues,
	}
}

func (n NPC) Name() string {
	return n.name
}
