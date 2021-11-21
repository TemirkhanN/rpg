package ui

import (
	"fmt"
	"unicode/utf8"

	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/internal/game/data"
	"github.com/TemirkhanN/rpg/pkg/rpg"
)

func DrawStatusPanel(player data.Player, on tcell.Screen) {
	NewBox(1, 1, 25, 10, "Status").Draw(on)
	NewText("Name: "+player.Name(), 2, 3).Draw(on, InfoTextStyle)
	NewText("Location: "+player.Whereabouts().Name(), 2, 4).Draw(on, InfoTextStyle)
}

func DrawPlayer(player data.Player, on tcell.Screen) {
	// Draw player.
	DrawUnit(player, on)

	// Draw player dialogue panel.
	dialogueDefaultPosition := NewPosition(80, 1)
	DrawDialogue(player.ActiveDialogue(), on, dialogueDefaultPosition)
}

func DrawDialogue(dialogue rpg.Dialogue, on tcell.Screen, at Pos) {
	if dialogue.Empty() {
		return
	}

	charactersPerLine := 26
	textStartAt := NewPosition(at.X()+2, at.Y()+1)

	textEndsAt := drawTextWithAutoLinebreaks(on, string(dialogue.Text()), textStartAt, charactersPerLine)

	if len(dialogue.Choices()) != 0 {
		textEndsAt.y++
	}

	for i, reply := range dialogue.Choices() {
		textEndsAt = drawTextWithAutoLinebreaks(
			on,
			fmt.Sprintf("%d.%s", i+1, reply),
			textEndsAt,
			charactersPerLine,
		)
	}

	NewBox(at.X(), at.Y(), at.X()+charactersPerLine+4, textEndsAt.Y(), dialogue.With().Name()).Draw(on, BoxStyle)
}

func drawTextWithAutoLinebreaks(on tcell.Screen, text string, at Pos, charactersPerLine int) Pos {
	verticalOffset := at.Y()
	dialogueLength := utf8.RuneCountInString(text)
	if dialogueLength <= charactersPerLine {
		NewText(text, 82, verticalOffset).Draw(on, TextStyle)
		verticalOffset++
	} else {
		var linedString string
		for i, r := range text {
			linedString += string(r)

			if (i+1)%charactersPerLine == 0 || i+1 == dialogueLength {
				NewText(linedString, 82, verticalOffset).Draw(on, TextStyle)
				verticalOffset++
				linedString = ""
			}
		}
	}

	return NewPosition(at.X(), verticalOffset)
}
