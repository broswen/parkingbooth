package ticket

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo TicketRepository
}

func NewService(repo TicketRepository) (*Service, error) {
	return &Service{
		repo: repo,
	}, nil
}

func (ts *Service) GenerateTicket(location string) (Ticket, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Ticket{}, err
	}

	t := Ticket{
		Id:       id.String(),
		Location: location,
		Start:    time.Now().Unix(),
	}
	t, err = ts.repo.SaveTicket(t)
	if err != nil {
		return Ticket{}, err
	}

	return t, nil
}

func (ts *Service) GetTicket(location, id string) (Ticket, error) {
	t, err := ts.repo.GetTicket(location, id)
	if err != nil {
		return Ticket{}, err
	}
	return t, nil
}

func (ts *Service) CompleteTicket(location, id string) (Ticket, error) {
	ticket, err := ts.GetTicket(location, id)
	if err != nil {
		return Ticket{}, err
	}
	ticket.Stop = time.Now().Unix()
	ticket.Duration = ticket.Stop - ticket.Start
	ticket, err = ts.repo.SaveTicket(ticket)
	if err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}

func (ts *Service) PayTicket(location, id string, payment string) (Ticket, error) {
	ticket, err := ts.GetTicket(location, id)
	if err != nil {
		return Ticket{}, err
	}

	ticket.Payment = payment
	ticket, err = ts.repo.SaveTicket(ticket)
	if err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}
