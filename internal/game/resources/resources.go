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

type npc struct {
	ID   int
	Name string
	Icon UnicodeIcon
}

type UnicodeIcon rune

func (r *UnicodeIcon) UnmarshalYAML(n *yaml.Node) error {
	var s string
	if err := n.Decode(&s); err != nil {
		return errors.Unwrap(err)
	}

	rn, _ := utf8.DecodeRune([]byte(s))
	*r = UnicodeIcon(rn)

	return nil
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
	ID   int
	Name string
	Npcs []struct {
		ID       int
		Position position
	}
}

//go:embed *
var resourceDir embed.FS

type Resources struct {
	initialized bool
	npcs        map[int]data.Npc
}

func LoadResources() Resources {
	res := Resources{
		initialized: true,
		npcs:        nil,
	}

	res.loadNpcs()

	return res
}

func (r *Resources) loadNpcs() {
	npcFile, err := resourceDir.ReadFile("npc.yaml")
	if err != nil {
		panic(err)
	}

	allNpc := make([]npc, 0)
	_ = yaml.Unmarshal(npcFile, &allNpc)

	r.npcs = make(map[int]data.Npc, len(allNpc))

	for _, npcDetails := range allNpc {
		r.npcs[npcDetails.ID] = *data.NewNpc(
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

func (r Resources) GetNPC(id int) (data.Npc, error) {
	npc, exists := r.npcs[id]
	if !exists {
		return data.Npc{}, fmt.Errorf("NPC with id %d does not exist", id)
	}

	return npc, nil
}

func (r Resources) LoadLocation(id int) (data.Location, error) {
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

		loc := data.NewLocation(rpg.NewLocation(locationDetails.ID, locationDetails.Name))

		for _, npcOnLocation := range locationDetails.Npcs {
			npc, err := r.GetNPC(npcOnLocation.ID)
			if err != nil {
				panic(err)
			}

			loc.Spawn(npc, data.NewPos(npcOnLocation.Position.X, npcOnLocation.Position.Y))
		}

		return loc, nil
	}

	return data.Location{}, fmt.Errorf("location with id %d does not exist", id)
}
