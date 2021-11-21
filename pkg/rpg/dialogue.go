package rpg

type Phrase string

type Dialogue struct {
	with    NPC
	text    Phrase
	choices []Phrase
}

func NewDialogue(text Phrase, choices []Phrase) Dialogue {
	return Dialogue{
		with:    NoNpc,
		text:    text,
		choices: choices,
	}
}

func (d *Dialogue) With() NPC {
	return d.with
}

func (d Dialogue) Text() Phrase {
	return d.text
}

func (d Dialogue) Choices() []Phrase {
	return d.choices
}

func (d Dialogue) Empty() bool {
	return d.Text() == ""
}

var NoDialogue = NewDialogue("", nil)
