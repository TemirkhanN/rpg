package rpg

type NPC struct {
	name string
}

func NewNPC(NPCName string) NPC {
	return NPC{
		name: NPCName,
	}
}

func (n NPC) Name() string {
	return n.name
}
