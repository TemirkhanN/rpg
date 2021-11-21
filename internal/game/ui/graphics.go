package ui

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"

	"github.com/TemirkhanN/rpg/internal/game/data"
)

type Drawer interface {
	Icon() rune
	Style() tcell.Style
	Position() data.Position
}

type Pos struct {
	x int
	y int
}

func NewPosition(x int, y int) Pos {
	return Pos{
		x: x,
		y: y,
	}
}

func (p Pos) X() int {
	return p.x
}

func (p Pos) Y() int {
	return p.y
}

type Box struct {
	title       *Text
	leftTop     Pos
	rightBottom Pos
}

type Text struct {
	text     string
	position Pos
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
		leftTop:     Pos{x: x1, y: y1},
		rightBottom: Pos{x: x2, y: y2},
	}
}

func NewText(value string, posX int, posY int) *Text {
	return &Text{
		text:     value,
		position: Pos{x: posX, y: posY},
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

	horizontalOffset := t.position.X()
	for _, c := range t.text {
		var comb []rune

		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}

		on.SetContent(horizontalOffset, t.position.Y(), c, comb, style)
		horizontalOffset += w
	}
}

func (b Box) Draw(on tcell.Screen, useStyle ...tcell.Style) {
	style := BoxStyle
	if len(useStyle) == 1 {
		style = useStyle[0]
	}

	x1 := b.leftTop.X()
	x2 := b.rightBottom.X()
	y1 := b.leftTop.Y()
	y2 := b.rightBottom.Y()

	hasTitle := false
	titleStartAt := b.title.position.X()
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

func DrawUnit(unit Drawer, on tcell.Screen) {
	NewText(string(unit.Icon()), unit.Position().X(), unit.Position().Y()).Draw(on, unit.Style())
}

func DrawLocation(location data.Location, on tcell.Screen) {
	NewBox(
		location.LeftTop().X(),
		location.LeftTop().Y(),
		location.RightBottom().X(),
		location.RightBottom().Y(),
		location.Name(),
	).Draw(on, BoxStyle)

	for _, gameNpc := range location.Npcs() {
		DrawUnit(gameNpc, on)
	}

	for _, passage := range location.Passages() {
		DrawUnit(passage, on)
	}

	for _, object := range location.Objects() {
		DrawUnit(object, on)
	}
}

var (
	TextStyle     = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	InfoTextStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorCadetBlue)
	BoxStyle      = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
)
