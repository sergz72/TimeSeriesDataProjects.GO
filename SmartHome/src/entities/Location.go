package entities

type Location struct {
	Id           int
	Name         string
	LocationType string
}

func (l Location) GetId() int {
	return l.Id
}
