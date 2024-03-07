package entities

type Category struct {
	Id   int
	Name string
}

func (c Category) GetId() int {
	return c.Id
}
