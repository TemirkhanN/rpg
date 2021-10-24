package rpg

type Location struct {
	name string
}

func NewLocation(locationName string) Location {
	return Location{
		name: locationName,
	}
}

func (l Location) Name() string {
	return l.name
}
