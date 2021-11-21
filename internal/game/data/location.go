package data

import (
	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type Location interface {
	AddPassage(in Position, out Position, to Location)
	LeftTop() Position
	RightBottom() Position
	Name() string
	Spawn(npc Npc, position Position)
	Npcs() []*Npc
	Objects() []*Object
	Passages() []Passage
}

type CommonLocation struct {
	leftTop     Position
	rightBottom Position
	whereabouts rpg.Location
	npcs        []*Npc
	objects     []*Object
	passages    []Passage
}

func NewLocation(whereabouts rpg.Location) CommonLocation {
	return CommonLocation{
		leftTop:     NewPos(26, 1),
		rightBottom: NewPos(79, 20),
		whereabouts: whereabouts,
		npcs:        nil,
		objects:     nil,
		passages:    nil,
	}
}

func (l *CommonLocation) AddPassage(in Position, out Position, to Location) {
	passage := Passage{
		in:  in,
		to:  to,
		out: out,
	}

	l.passages = append(l.passages, passage)
}

func (l CommonLocation) LeftTop() Position {
	return l.leftTop
}

func (l CommonLocation) RightBottom() Position {
	return l.rightBottom
}

func (l CommonLocation) Name() string {
	return l.whereabouts.Name()
}

func (l *CommonLocation) Spawn(npc Npc, position Position) {
	spawnedNpc := npc
	spawnedNpc.position = position

	l.npcs = append(l.npcs, &spawnedNpc)
}

func (l CommonLocation) Npcs() []*Npc {
	return l.npcs
}

func (l CommonLocation) Objects() []*Object {
	return l.objects
}

func (l CommonLocation) Passages() []Passage {
	return l.passages
}

func positionInsideLocation(pos Position, location Location) bool {
	if location.LeftTop().X() >= pos.X() || location.LeftTop().Y() >= pos.Y() {
		return false
	}

	if location.RightBottom().X()-1 <= pos.X() || location.RightBottom().Y() <= pos.Y() {
		return false
	}

	return true
}
