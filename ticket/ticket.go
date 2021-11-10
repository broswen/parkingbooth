package ticket

import (
	"time"

	"github.com/broswen/parkingbooth/repository"
	"github.com/google/uuid"
)

type Service struct {
	repo repository.TicketRepository
}

func NewService(repo repository.TicketRepository) (*Service, error) {
	return &Service{
		repo: repo,
	}, nil
}

func (ts *Service) GenerateTicket(location string) (repository.Ticket, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return repository.Ticket{}, err
	}

	t := repository.Ticket{
		Id:       id.String(),
		Location: location,
		Start:    time.Now().Unix(),
	}
	t, err = ts.repo.SaveTicket(t)
	if err != nil {
		return repository.Ticket{}, err
	}

	return t, nil
}

func (ts *Service) GetTicket(location, id string) (repository.Ticket, error) {
	t, err := ts.repo.GetTicket(location, id)
	if err != nil {
		return repository.Ticket{}, err
	}
	return t, nil
}

func (ts *Service) CompleteTicket(location, id string) (repository.Ticket, error) {
	ticket, err := ts.GetTicket(location, id)
	if err != nil {
		return repository.Ticket{}, err
	}
	ticket.Stop = time.Now().Unix()
	ticket.Duration = ticket.Stop - ticket.Start
	ticket, err = ts.repo.SaveTicket(ticket)
	if err != nil {
		return repository.Ticket{}, err
	}

	return ticket, nil
}

func (ts *Service) PayTicket(location, id string, payment string) (repository.Ticket, error) {
	ticket, err := ts.GetTicket(location, id)
	if err != nil {
		return repository.Ticket{}, err
	}

	ticket.Payment = payment
	ticket, err = ts.repo.SaveTicket(ticket)
	if err != nil {
		return repository.Ticket{}, err
	}

	return ticket, nil
}
