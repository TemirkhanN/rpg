package rpg

type Dialogue struct {
	text    string
	choices []string
}

func NewDialogue(text string, choices []string) Dialogue {
	return Dialogue{
		text:    text,
		choices: choices,
	}
}

func (d Dialogue) Text() string {
	return d.text
}

func (d Dialogue) Choices() []string {
	return d.choices
}
