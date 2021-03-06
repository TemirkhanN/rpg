package resources

import (
	"embed"
	"errors"
	"fmt"
	"unicode/utf8"

	"gopkg.in/yaml.v3"

	"github.com/TemirkhanN/rpg/internal/game/data"
	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type unicodeIcon rune

type id int

func (r *unicodeIcon) UnmarshalYAML(n *yaml.Node) error {
	var s string
	if err := n.Decode(&s); err != nil {
		return errors.Unwrap(err)
	}

	rn, _ := utf8.DecodeRune([]byte(s))
	*r = unicodeIcon(rn)

	return nil
}

type object struct {
	ID    id
	Icon  unicodeIcon
	Solid bool
	Name  string
}

type npc struct {
	ID   id
	Name string
	Icon unicodeIcon
}

type dialogue struct {
	Text    string
	Replies []string
}

type position struct {
	X int
	Y int
}

type location struct {
	ID   id
	Name string
	Npcs []struct {
		ID       id
		Position position
	}
	Passages []struct {
		In  position
		Out position
		To  id
	}
	Objects []struct {
		ID       id
		Position position
	}
}

//go:embed *
var resourceDir embed.FS

type Resources struct {
	initialized bool
	npcs        map[id]data.Npc
	objects     map[id]data.Object
	locations   map[id]data.Location
}

func LoadResources() Resources {
	res := Resources{
		initialized: true,
		npcs:        nil,
		objects:     nil,
		locations:   make(map[id]data.Location),
	}

	res.loadObjects()
	res.loadNpcs()

	return res
}

func (r *Resources) loadObjects() {
	yamlFile, err := resourceDir.ReadFile("object.yaml")
	if err != nil {
		panic(err)
	}

	allObjects := make([]object, 0)
	_ = yaml.Unmarshal(yamlFile, &allObjects)

	r.objects = make(map[id]data.Object, len(allObjects))

	for _, objectDetails := range allObjects {
		r.objects[objectDetails.ID] = data.NewObject(
			objectDetails.Name,
			rune(objectDetails.Icon),
			objectDetails.Solid,
		)
	}
}

func (r Resources) GetObject(id id) (data.Object, error) {
	object, exists := r.objects[id]
	if !exists {
		return data.Object{}, fmt.Errorf("object with id %d does not exist", id)
	}

	return object, nil
}

func (r *Resources) loadNpcs() {
	yamlFile, err := resourceDir.ReadFile("npc.yaml")
	if err != nil {
		panic(err)
	}

	allNpc := make([]npc, 0)
	_ = yaml.Unmarshal(yamlFile, &allNpc)

	r.npcs = make(map[id]data.Npc, len(allNpc))

	for _, npcDetails := range allNpc {
		r.npcs[npcDetails.ID] = data.NewNpc(
			rpg.NewNPC(npcDetails.Name, loadDialogues(npcDetails)),
			rune(npcDetails.Icon),
			data.NoPosition,
		)
	}
}

func loadDialogues(n npc) map[rpg.Phrase]rpg.Dialogue {
	yamlFile, err := resourceDir.ReadFile(fmt.Sprintf("dialogue/%s.yaml", n.Name))
	if err != nil {
		return nil
	}

	npcDialogues := make(map[string]dialogue)
	_ = yaml.Unmarshal(yamlFile, &npcDialogues)

	dialogues := make(map[rpg.Phrase]rpg.Dialogue, len(npcDialogues))
	for reply, dialogue := range npcDialogues {
		phrase := rpg.Phrase(dialogue.Text)

		replies := make([]rpg.Phrase, len(dialogue.Replies))
		for i, parsedReply := range dialogue.Replies {
			replies[i] = rpg.Phrase(parsedReply)
		}

		dialogues[rpg.Phrase(reply)] = rpg.NewDialogue(phrase, replies)
	}

	return dialogues
}

func (r Resources) GetNPC(id id) (data.Npc, error) {
	npc, exists := r.npcs[id]
	if !exists {
		return data.Npc{}, fmt.Errorf("NPC with id %d does not exist", id)
	}

	return npc, nil
}

func (r Resources) LoadLocation(id id) (nilLoc data.Location, err error) {
	loc, exists := r.locations[id]
	if exists {
		return loc, nil
	}

	yamlFile, err := resourceDir.ReadFile("location.yaml")
	if err != nil {
		panic(err)
	}

	var allLocations []location
	_ = yaml.Unmarshal(yamlFile, &allLocations)

	for _, locationDetails := range allLocations {
		if locationDetails.ID != id {
			continue
		}

		loc := data.NewLocation(rpg.NewLocation(int(locationDetails.ID), locationDetails.Name))

		for _, passageDetails := range locationDetails.Passages {
			leadsToID := passageDetails.To

			leadsTo := lazyLoadLocation{
				location: nil,
				loader: func() data.Location {
					leadsTo, err := r.LoadLocation(leadsToID)
					if err != nil {
						panic(err)
					}

					return leadsTo
				},
			}

			if err != nil {
				return nilLoc, fmt.Errorf("location %d has passage to unknown location %d", id, leadsToID)
			}

			loc.AddPassage(
				data.NewPos(passageDetails.In.X, passageDetails.In.Y),
				data.NewPos(passageDetails.Out.X, passageDetails.Out.Y),
				&leadsTo,
			)
		}

		for _, npcOnLocation := range locationDetails.Npcs {
			npc, err := r.GetNPC(npcOnLocation.ID)
			if err != nil {
				panic(err)
			}

			loc.Spawn(npc, data.NewPos(npcOnLocation.Position.X, npcOnLocation.Position.Y))
		}

		for _, objectInLocation := range locationDetails.Objects {
			object, _ := r.GetObject(objectInLocation.ID)

			loc.PlaceObject(object, data.NewPos(objectInLocation.Position.X, objectInLocation.Position.Y))
		}

		r.locations[id] = &loc

		return r.locations[id], nil
	}

	return nilLoc, fmt.Errorf("location with id %d does not exist", id)
}

type lazyLoadLocation struct {
	location data.Location
	loader   func() data.Location
}

func (ll *lazyLoadLocation) load() {
	if ll.location == nil {
		ll.location = ll.loader()
	}
}

func (ll *lazyLoadLocation) AddPassage(in data.Position, out data.Position, to data.Location) {
	ll.load()

	ll.location.AddPassage(in, out, to)
}

func (ll lazyLoadLocation) LeftTop() data.Position {
	ll.load()

	return ll.location.LeftTop()
}

func (ll lazyLoadLocation) RightBottom() data.Position {
	ll.load()

	return ll.location.RightBottom()
}

func (ll lazyLoadLocation) Name() string {
	ll.load()

	return ll.location.Name()
}

func (ll *lazyLoadLocation) Spawn(npc data.Npc, position data.Position) {
	ll.load()

	ll.location.Spawn(npc, position)
}

func (ll lazyLoadLocation) Npcs() []data.Npc {
	ll.load()

	return ll.location.Npcs()
}

func (ll *lazyLoadLocation) PlaceObject(object data.Object, at data.Position) {
	ll.load()

	ll.location.PlaceObject(object, at)
}

func (ll lazyLoadLocation) Objects() []data.Object {
	ll.load()

	return ll.location.Objects()
}

func (ll lazyLoadLocation) Passages() []data.Passage {
	ll.load()

	return ll.location.Passages()
}
