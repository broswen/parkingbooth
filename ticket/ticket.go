package ticket

import (
	"time"

	"github.com/broswen/parkingbooth/location"
	"github.com/google/uuid"
)

type Service struct {
	ticketRepo   TicketRepository
	locationRepo location.LocationRepository
}

func NewService(ticketRepo TicketRepository, locationRepo location.LocationRepository) (*Service, error) {
	return &Service{
		ticketRepo:   ticketRepo,
		locationRepo: locationRepo,
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
	t, err = ts.ticketRepo.SaveTicket(t)
	if err != nil {
		return Ticket{}, err
	}

	return t, nil
}

func (ts *Service) GetTicket(location, id string) (Ticket, error) {
	t, err := ts.ticketRepo.GetTicket(location, id)
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
	ticket, err = ts.ticketRepo.SaveTicket(ticket)
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
	ticket, err = ts.ticketRepo.SaveTicket(ticket)
	if err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}
