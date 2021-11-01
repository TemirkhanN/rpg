package resources

import (
	"embed"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/TemirkhanN/rpg/pkg/rpg"
)

type npc struct {
	Name string
}

type dialogue struct {
	Text    string
	Replies []string
}

//go:embed *
var resourceDir embed.FS

type Resources struct {
	initialized bool
	npcs        map[string]rpg.NPC
	locations   map[string]rpg.Location
}

func LoadResources() Resources {
	res := Resources{
		initialized: true,
		npcs:        nil,
		locations:   nil,
	}
	rd := resourceDir
	npcFile, err := rd.ReadFile("npc.yaml")
	if err != nil {
		panic(err)
	}

	allNpc := make([]npc, 0)
	_ = yaml.Unmarshal(npcFile, &allNpc)

	res.npcs = make(map[string]rpg.NPC, len(allNpc))

	for _, n := range allNpc {
		res.npcs[n.Name] = rpg.NewNPC(n.Name, n.loadDialogues())
	}

	locFile, err := rd.ReadFile("location.yaml")
	if err != nil {
		panic(err)
	}

	allLocations := make(map[string]interface{})
	_ = yaml.Unmarshal(locFile, allLocations)

	res.locations = make(map[string]rpg.Location, len(allLocations))

	for locationName := range allLocations {
		res.locations[locationName] = rpg.NewLocation(locationName)
	}

	return res
}

func (n npc) loadDialogues() map[string]rpg.Dialogue {
	yamlFile, err := resourceDir.ReadFile(fmt.Sprintf("dialogue/%s.yaml", n.Name))
	if err != nil {
		return nil
	}

	npcDialogues := make(map[string]dialogue)
	_ = yaml.Unmarshal(yamlFile, &npcDialogues)

	dialogues := make(map[string]rpg.Dialogue, len(npcDialogues))
	for reply, dialogue := range npcDialogues {
		dialogues[reply] = rpg.NewDialogue(dialogue.Text, dialogue.Replies)
	}

	return dialogues
}

func (r Resources) GetNPC(name string) (rpg.NPC, error) {
	npc, exists := r.npcs[name]
	if !exists {
		return rpg.NPC{}, fmt.Errorf("NPC with name %s does not exist", name)
	}

	return npc, nil
}

func (r Resources) Locations() []rpg.Location {
	locations := make([]rpg.Location, 0, len(r.locations))

	for _, value := range r.locations {
		locations = append(locations, value)
	}

	return locations
}

func (r Resources) GetLocation(name string) (rpg.Location, error) {
	location, exists := r.locations[name]
	if !exists {
		return rpg.Location{}, fmt.Errorf("location with name %s does not exist", name)
	}

	return location, nil
}
