package rpg

type Location struct {
	id   int
	name string
}

func NewLocation(id int, locationName string) Location {
	return Location{
		id:   id,
		name: locationName,
	}
}

func (l Location) Name() string {
	return l.name
}
