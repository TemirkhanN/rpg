package data

import (
	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type Location struct {
	leftTop     Position
	rightBottom Position
	whereabouts rpg.Location
	npcs        []*Npc
	objects     []*Object
	passages    []Passage
}

func NewLocation(whereabouts rpg.Location) Location {
	return Location{
		leftTop:     NewPos(26, 1),
		rightBottom: NewPos(79, 20),
		whereabouts: whereabouts,
		npcs:        nil,
		objects:     nil,
		passages:    nil,
	}
}

func (l Location) LeftTop() Position {
	return l.leftTop
}

func (l Location) RightBottom() Position {
	return l.rightBottom
}

func positionInsideLocation(pos Position, location Location) bool {
	if location.leftTop.X() >= pos.X() || location.leftTop.Y() >= pos.Y() {
		return false
	}

	if location.rightBottom.X()-1 <= pos.X() || location.rightBottom.Y() <= pos.Y() {
		return false
	}

	return true
}

func (l Location) Name() string {
	return l.whereabouts.Name()
}

func (l *Location) Spawn(npc Npc, position Position) {
	spawnedNpc := npc
	spawnedNpc.position = position

	l.npcs = append(l.npcs, &spawnedNpc)
}

func (l Location) Npcs() []*Npc {
	return l.npcs
}

func (l Location) Objects() []*Object {
	return l.objects
}

func (l Location) Passages() []Passage {
	return l.passages
}
