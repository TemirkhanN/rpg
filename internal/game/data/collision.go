package data

type Collider interface {
	Collides(with Position) bool
}

func collides(p1 Position, p2 Position) bool {
	if p1.Y() != p2.Y() {
		return false
	}

	if p1.X() < p2.X()-1 || p1.X() > p2.X()+1 {
		return false
	}

	return true
}
