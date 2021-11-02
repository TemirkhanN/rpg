package ui

import (
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

type asci struct {
	symbol rune
	style  tcell.Style
}

type Position struct {
	X int
	Y int
}

type Box struct {
	title       *Text
	leftTop     Position
	rightBottom Position
}

type Text struct {
	text     string
	position Position
}

func (p Position) Inside(a Area) bool {
	if a.rectangle.leftTop.X >= p.X || a.rectangle.leftTop.Y >= p.Y {
		return false
	}

	if a.rectangle.rightBottom.X-1 <= p.X || a.rectangle.rightBottom.Y <= p.Y {
		return false
	}

	return true
}

func CalculateTextWidth(text string) int {
	width := 0
	for _, c := range text {
		width += runewidth.RuneWidth(c)
	}

	return width
}

func NewBox(x1 int, y1 int, x2 int, y2 int, title ...string) Box {
	if y2 < y1 {
		y1, y2 = y2, y1
	}

	if x2 < x1 {
		x1, x2 = x2, x1
	}

	boxTitle := NewText("", 0, 0)
	if len(title) == 1 {
		titleLen := len(title[0])
		titleStartAt := (x2-x1-titleLen)/2 + x1

		if titleStartAt < x1 {
			titleStartAt = x1
		}

		boxTitle = NewText(title[0], titleStartAt, y1)
	}

	return Box{
		title:       boxTitle,
		leftTop:     Position{X: x1, Y: y1},
		rightBottom: Position{X: x2, Y: y2},
	}
}

func NewText(value string, posX int, posY int) *Text {
	return &Text{
		text:     value,
		position: Position{X: posX, Y: posY},
	}
}

func CreateScreen() tcell.Screen {
	screen, e := tcell.NewScreen()

	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	if e := screen.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)

	screen.SetStyle(defStyle)

	return screen
}

func (t Text) Draw(on tcell.Screen, useStyle ...tcell.Style) {
	style := TextStyle
	if len(useStyle) == 1 {
		style = useStyle[0]
	}

	horizontalOffset := t.position.X
	for _, c := range t.text {
		var comb []rune

		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}

		on.SetContent(horizontalOffset, t.position.Y, c, comb, style)
		horizontalOffset += w
	}
}

func (b Box) Draw(on tcell.Screen, useStyle ...tcell.Style) {
	style := BoxStyle
	if len(useStyle) == 1 {
		style = useStyle[0]
	}

	x1 := b.leftTop.X
	x2 := b.rightBottom.X
	y1 := b.leftTop.Y
	y2 := b.rightBottom.Y

	hasTitle := false
	titleStartAt := b.title.position.X
	titleEndAt := titleStartAt + len(b.title.text) - 1
	if b.title.text != "" {
		hasTitle = true
		b.title.Draw(on, TextStyle)
	}

	for col := x1; col <= x2; col++ {
		if !hasTitle || col < titleStartAt-1 || col > titleEndAt+1 {
			on.SetContent(col, y1, tcell.RuneHLine, nil, style)
		}

		on.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}

	for row := y1 + 1; row < y2; row++ {
		on.SetContent(x1, row, tcell.RuneVLine, nil, style)
		on.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	if y1 != y2 && x1 != x2 {
		on.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		on.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		on.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		on.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
}

func DrawTextWithAutoLinebreaks(on tcell.Screen, text string, at Position, charactersPerLine int) Position {
	verticalOffset := at.Y
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

	return Position{X: at.X, Y: verticalOffset}
}

var (
	TextStyle     = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	InfoTextStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorCadetBlue)
	BoxStyle      = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
)
