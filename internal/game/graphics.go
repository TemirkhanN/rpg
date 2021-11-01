package game

import (
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"

	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type teleport struct {
	locations []rpg.Location
	panel     box
}

func newTeleport(locations []rpg.Location, panel box) teleport {
	return teleport{
		locations: locations,
		panel:     panel,
	}
}

func (t teleport) click(player *rpg.Player, pos position) {
	leftTop := t.panel.leftTop
	for i, location := range t.locations {
		optionPos := position{leftTop.x + 1, leftTop.y + i + 1}
		// if option is not on the same line or text is not within clicked position.
		if pos.y != optionPos.y ||
			pos.x < optionPos.x ||
			pos.x > optionPos.x+calculateTextWidth(location.Name()) {
			continue
		}

		if player.Whereabouts().Name() != location.Name() {
			player.MoveToLocation(location)
		}

		break
	}
}

func (t teleport) draw(on tcell.Screen, panelStyle tcell.Style, optionsStyle tcell.Style) {
	t.panel.draw(on, panelStyle)

	leftTop := t.panel.leftTop
	for i, location := range t.locations {
		optionPos := position{leftTop.x + 1, leftTop.y + i + 1}
		newText(fmt.Sprintf("* %s", location.Name()), optionPos.x, optionPos.y).draw(on, optionsStyle)
	}
}

type npc struct {
	asci asci
	npc  rpg.NPC
	pos  position
}

func (n npc) collides(with position) bool {
	if with.y != n.pos.y {
		return false
	}

	// npc takes 3 slots on horizontal line
	if with.x < n.pos.x-1 || with.x > n.pos.x+1 {
		return false
	}

	return true
}

type player struct {
	asci                asci
	player              *rpg.Player
	pos                 position
	currentDialogue     rpg.Dialogue
	currentDialogueWith rpg.NPC
}

type asci struct {
	symbol rune
	style  tcell.Style
}

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

func calculateTextWidth(text string) int {
	width := 0
	for _, c := range text {
		width += runewidth.RuneWidth(c)
	}

	return width
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

func drawTextWithAutoLinebreaks(on tcell.Screen, text string, at position, charactersPerLine int) position {
	verticalOffset := at.y
	dialogueLength := utf8.RuneCountInString(text)
	if dialogueLength <= charactersPerLine {
		newText(text, 82, verticalOffset).draw(on, textStyle)
		verticalOffset++
	} else {
		var linedString string
		for i, r := range text {
			linedString += string(r)

			if (i+1)%charactersPerLine == 0 || i+1 == dialogueLength {
				newText(linedString, 82, verticalOffset).draw(on, textStyle)
				verticalOffset++
				linedString = ""
			}
		}
	}

	return position{x: at.x, y: verticalOffset}
}

var (
	textStyle        = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	infoTextStyle    = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorCadetBlue)
	boxStyle         = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	friendlyNPCStyle = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorGreen).Bold(true)
	playerStyle      = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGoldenrod).Bold(true)
)
