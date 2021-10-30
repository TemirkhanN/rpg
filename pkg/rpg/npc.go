package rpg

type NPC struct {
	name string
}

func NewNPC(name string) NPC {
	return NPC{
		name: name,
	}
}

func (n NPC) Name() string {
	return n.name
}
