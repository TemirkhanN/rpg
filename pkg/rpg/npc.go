package rpg

type NPC struct {
	name      string
	dialogues map[Phrase]Dialogue
}

func NewNPC(name string, dialogues map[Phrase]Dialogue) NPC {
	return NPC{
		name:      name,
		dialogues: dialogues,
	}
}

func (n NPC) StartConversation(p Player) Dialogue {
	return n.Reply("defaultDialogue")
}

func (n NPC) Reply(on Phrase) Dialogue {
	dialogue := n.dialogues[on]
	dialogue.with = n

	return dialogue
}

func (n NPC) Name() string {
	return n.name
}

var NoNpc = NewNPC("", nil)
