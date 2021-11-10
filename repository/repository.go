package repository

import (
	"fmt"
	"time"
)

type Ticket struct {
	Id       string        `json:"id"`
	Location string        `json:"location"`
	Start    time.Time     `json:"start"`
	Stop     time.Time     `json:"stop,omitempty"`
	Duration time.Duration `json:"duration,omitempty"`
	Payment  string        `json:"payment,omitempty"`
}

type TicketRepository interface {
	GetTicket(location, id string) (Ticket, error)
	SaveTicket(t Ticket) (Ticket, error)
}

type MapRepository struct {
	m map[string]map[string]Ticket
}

func NewRepository(t string) (TicketRepository, error) {
	switch t {
	case "map":
		return MapRepository{
			m: make(map[string]map[string]Ticket, 0),
		}, nil
	case "postgres":
		return nil, fmt.Errorf("type not implemented: %v\n", t)
	default:
		return nil, fmt.Errorf("unknown repository type: %v\n", t)
	}
}

func (mr MapRepository) GetTicket(location, id string) (Ticket, error) {
	l, ok := mr.m[location]
	if !ok {
		return Ticket{}, fmt.Errorf("no tickets for location: %v\n", location)
	}
	t, ok := l[id]
	if !ok {
		return Ticket{}, fmt.Errorf("no tickets for id: %v\n", id)
	}
	return t, nil
}

func (mr MapRepository) SaveTicket(t Ticket) (Ticket, error) {
	l, ok := mr.m[t.Location]
	if !ok {
		mr.m[t.Location] = make(map[string]Ticket, 0)
	}
	l = mr.m[t.Location]
	l[t.Id] = t
	return t, nil
}
