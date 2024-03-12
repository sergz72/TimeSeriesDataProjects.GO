package entities

type Sensor struct {
	Id            int
	Name          string
	DataType      string
	LocationId    int
	DeviceId      int
	DeviceSensors map[int]string
	Offsets       map[string]int
}

func (s Sensor) GetId() int {
	return s.Id
}
