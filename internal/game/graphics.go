package game

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

type position struct {
	x int
	y int
}

type box struct {
	title       *text
	leftTop     position
	rightBottom position
}

type text struct {
	text     string
	position position
}

func newBox(x1 int, y1 int, x2 int, y2 int, title ...string) *box {
	if y2 < y1 {
		y1, y2 = y2, y1
	}

	if x2 < x1 {
		x1, x2 = x2, x1
	}

	boxTitle := newText("", 0, 0)
	if len(title) == 1 {
		titleLen := len(title[0])
		titleStartAt := (x2-x1-titleLen)/2 + x1

		if titleStartAt < x1 {
			titleStartAt = x1
		}

		boxTitle = newText(title[0], titleStartAt, y1)
	}

	return &box{
		title: boxTitle,
		leftTop: position{
			x: x1,
			y: y1,
		},
		rightBottom: position{
			x: x2,
			y: y2,
		},
	}
}

func newText(value string, posX int, posY int) *text {
	return &text{
		text: value,
		position: position{
			x: posX,
			y: posY,
		},
	}
}

func createScreen() tcell.Screen {
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

func (t text) draw(on tcell.Screen, style tcell.Style) {
	horizontalOffset := t.position.x
	for _, c := range t.text {
		var comb []rune

		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}

		on.SetContent(horizontalOffset, t.position.y, c, comb, style)
		horizontalOffset += w
	}
}

func (b box) draw(on tcell.Screen, style tcell.Style) {
	x1 := b.leftTop.x
	x2 := b.rightBottom.x
	y1 := b.leftTop.y
	y2 := b.rightBottom.y

	hasTitle := false
	titleStartAt := b.title.position.x
	titleEndAt := titleStartAt + len(b.title.text) - 1
	if b.title.text != "" {
		hasTitle = true
		b.title.draw(on, textStyle)
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

var (
	textStyle     = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	infoTextStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorCadetBlue)
	boxStyle      = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
)
