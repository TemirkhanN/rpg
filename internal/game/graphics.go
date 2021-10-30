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
	title       string
	leftTop     position
	rightBottom position
}

func newBox(x1 int, y1 int, x2 int, y2 int, title ...string) box {
	if y2 < y1 {
		y1, y2 = y2, y1
	}

	if x2 < x1 {
		x1, x2 = x2, x1
	}

	var boxTitle string
	if len(title) == 1 {
		boxTitle = title[0]
	}

	return box{
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

func drawText(s tcell.Screen, x int, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune

		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}

		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func (b box) draw(s tcell.Screen, style tcell.Style) {
	x1 := b.leftTop.x
	x2 := b.rightBottom.x
	y1 := b.leftTop.y
	y2 := b.rightBottom.y

	var titleStartAt int
	var titleEndAt int

	if b.title != "" {
		titleLen := len(b.title)
		titleStartAt = (x2-x1-titleLen)/2 + x1
		titleEndAt = titleStartAt + titleLen

		if titleStartAt < x1 {
			titleStartAt = x1
		}

		drawText(s, titleStartAt, y1, textStyle, b.title)
	}

	hasTitle := titleStartAt != 0 && titleEndAt != 0

	for col := x1; col <= x2; col++ {
		if !hasTitle || col < titleStartAt || col > titleEndAt {
			s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		}

		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}

	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
}

var (
	textStyle     = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	infoTextStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorCadetBlue)
	boxStyle      = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
)
