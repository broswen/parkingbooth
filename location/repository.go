package location

import "fmt"

type Location struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type LocationRepository interface {
	GetLocation(id string) (Location, error)
	SaveLocation(l Location) (Location, error)
}

func NewMap() (LocationRepository, error) {
	return MapRepository{
		m: make(map[string]Location, 0),
	}, nil
}

type MapRepository struct {
	m map[string]Location
}

func (mr MapRepository) GetLocation(id string) (Location, error) {
	l, ok := mr.m[id]
	if !ok {
		return Location{}, fmt.Errorf("location not found")
	}
	return l, nil
}

func (mr MapRepository) SaveLocation(l Location) (Location, error) {
	mr.m[l.Id] = l
	return l, nil
}
