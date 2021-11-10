package ticket

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	Id       string        `json:"id"`
	Location string        `json:"location"`
	Start    time.Time     `json:"start"`
	Stop     time.Time     `json:"stop,omitempty"`
	Duration time.Duration `json:"duration,omitempty"`
	Payment  string        `json:"payment,omitempty"`
}

type Service struct{}

func NewService() (*Service, error) {
	return &Service{}, nil
}

func (ts *Service) GenerateTicket(location string) (Ticket, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Ticket{}, err
	}

	ticket := Ticket{
		Id:       id.String(),
		Location: location,
		Start:    time.Now(),
	}
	// TODO save ticket to database

	return ticket, nil
}

func (ts *Service) GetTicket(location, id string) (Ticket, error) {
	return Ticket{
		Id:       id,
		Location: location,
		Start:    time.Now(),
	}, nil
}

func (ts *Service) CompleteTicket(location, id string) (Ticket, error) {
	ticket, err := ts.GetTicket(location, id)
	if err != nil {
		return Ticket{}, err
	}
	ticket.Stop = time.Now()
	ticket.Duration = ticket.Stop.Sub(ticket.Start)
	// TODO save to database

	return ticket, nil
}

func (ts *Service) PayTicket(location, id string, payment string) (Ticket, error) {
	ticket, err := ts.GetTicket(location, id)
	if err != nil {
		return Ticket{}, err
	}

	ticket.Payment = payment
	// TODO save to database

	return ticket, nil
}
