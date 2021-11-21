package data

import (
	"github.com/gdamore/tcell"

	"github.com/TemirkhanN/rpg/internal/game/ui/icon"
	"github.com/TemirkhanN/rpg/internal/game/ui/interactive"
	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type Player struct {
	player   *rpg.Player
	location *Location
	position Position
}

func NewPlayer(player *rpg.Player) *Player {
	return &Player{
		player:   player,
		location: nil,
		position: NewPos(0, 0),
	}
}

func (p Player) ChooseDialogueReplyOption(choice int) {
	currentDialogue := p.player.ActiveDialogue()
	if currentDialogue.Empty() {
		return
	}

	availableReplies := currentDialogue.Choices()
	if len(availableReplies) == 0 {
		return
	}

	if choice < 0 || len(availableReplies) < choice {
		return
	}

	p.player.Reply(currentDialogue.With(), availableReplies[choice-1])
}

func (p Player) Name() string {
	return p.player.Name()
}

func (p Player) Whereabouts() Location {
	return *p.location
}

func (p *Player) StartDialogue(with Npc) {
	p.player.StartConversation(with.template)
}

func (p *Player) EndDialogue() {
	p.player.EndConversation()
}

func (p Player) ActiveDialogue() rpg.Dialogue {
	return p.player.ActiveDialogue()
}

func (p Player) Position() Position {
	return p.position
}

func (p *Player) MoveTo(to Position) {
	// position is horizontally too far from current player position
	if to.X() != p.position.X() && p.position.X()+1 != to.X() && p.position.X()-1 != to.X() {
		return
	}

	// position is vertically too far from current player position
	if to.Y() != p.position.Y() && p.position.Y()+1 != to.Y() && p.position.Y()-1 != to.Y() {
		return
	}

	currentLocation := p.location

	if !positionInsideLocation(to, *currentLocation) {
		return
	}

	for _, passage := range currentLocation.Passages() {
		if collides(to, passage.in) {
			p.TeleportToLocation(passage.to, passage.out)

			return
		}
	}

	for _, object := range currentLocation.Objects() {
		if object.Collides(to) {
			return
		}
	}

	for _, npc := range currentLocation.Npcs() {
		if npc.Collides(to) {
			p.StartDialogue(*npc)

			return
		}
	}

	if !p.ActiveDialogue().Empty() {
		p.EndDialogue()
	}

	p.position = to
}

func (p *Player) Teleport(to Position) {
	p.position = to
}

func (p *Player) TeleportToLocation(location *Location, to Position) {
	p.location = location
	p.Teleport(to)
}

func (p Player) Interact(element interactive.Element) {
	element.RunAction()
}

func (p Player) Icon() rune {
	return icon.DogFace
}

func (p Player) Style() tcell.Style {
	return playerIconStyle
}

func (p Player) Collides(with Position) bool {
	return collides(p.Position(), with)
}

var playerIconStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGoldenrod).Bold(true)
