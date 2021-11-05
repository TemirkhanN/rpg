package ui

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type Player struct {
	asci                asci
	player              *rpg.Player
	pos                 Position
	currentDialogue     rpg.Dialogue
	currentDialogueWith rpg.NPC
}

func NewPlayer(player *rpg.Player, pos Position, icon rune, iconStyle tcell.Style) Player {
	return Player{
		asci: asci{
			symbol: icon,
			style:  iconStyle,
		},
		currentDialogue:     rpg.NoDialogue,
		currentDialogueWith: rpg.NoNpc,
		player:              player,
		pos:                 pos,
	}
}

func (p Player) Player() *rpg.Player {
	return p.player
}

func (p Player) Draw(on tcell.Screen) {
	NewText(string(p.asci.symbol), p.pos.X, p.pos.Y).Draw(on, p.asci.style)
}

func (p Player) DrawDialogue(on tcell.Screen, at Position) {
	if p.currentDialogue.Empty() {
		return
	}

	charactersPerLine := 26
	textStartAt := Position{X: at.X + 2, Y: at.Y + 1}

	textEndsAt := DrawTextWithAutoLinebreaks(on, p.currentDialogue.Text(), textStartAt, charactersPerLine)

	if len(p.currentDialogue.Choices()) != 0 {
		textEndsAt.Y++
	}

	for i, reply := range p.currentDialogue.Choices() {
		textEndsAt = DrawTextWithAutoLinebreaks(
			on,
			fmt.Sprintf("%d.%s", i+1, reply),
			textEndsAt,
			charactersPerLine,
		)
	}

	NewBox(at.X, at.Y, at.X+charactersPerLine+4, textEndsAt.Y, p.currentDialogueWith.Name()).Draw(on, BoxStyle)
}

func (p *Player) ChooseDialogueReplyOption(input tcell.EventKey) {
	currentDialogue := p.currentDialogue
	if currentDialogue.Empty() {
		return
	}

	availableReplies := currentDialogue.Choices()
	if len(availableReplies) == 0 {
		return
	}

	choice, err := strconv.Atoi(string(input.Rune()))
	if err != nil || choice < 0 || len(availableReplies) < choice {
		return
	}

	dialogue := p.player.Reply(p.currentDialogueWith, availableReplies[choice-1])
	if dialogue.Text() == p.currentDialogue.Text() {
		return
	}

	p.currentDialogue = dialogue
}

func (p *Player) StartDialogue(npc NPC) {
	dialogue := p.player.StartConversation(*npc.npc)

	if dialogue.Text() == p.currentDialogue.Text() {
		return
	}

	p.currentDialogue = dialogue
	p.currentDialogueWith = *npc.npc
}

func (p *Player) EndDialogue() {
	if p.currentDialogue.Empty() {
		return
	}

	p.currentDialogueWith = rpg.NoNpc
	p.currentDialogue = rpg.NoDialogue
}

func (p *Player) MoveTo(pos Position) {
	if pos.X != p.pos.X && p.pos.X+1 != pos.X && p.pos.X-1 != pos.X {
		return
	}

	if pos.Y != p.pos.Y && p.pos.Y+1 != pos.Y && p.pos.Y-1 != pos.Y {
		return
	}

	p.pos = pos
}

func (p Player) Position() Position {
	return p.pos
}

var PlayerIconStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGoldenrod).Bold(true)
