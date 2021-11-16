package account

import (
	"time"

	"github.com/broswen/parkingbooth/location"
	"github.com/google/uuid"
)

type Service struct {
	accountRepo  AccountRepository
	locationRepo location.LocationRepository
}

func NewService(accountRepo AccountRepository, locationRepo location.LocationRepository) (*Service, error) {
	return &Service{
		accountRepo:  accountRepo,
		locationRepo: locationRepo,
	}, nil
}

func (as *Service) CreateAccount(a Account) (Account, error) {
	return as.accountRepo.CreateAccount(a)
}

func (as *Service) UpdateAccount(a Account) (Account, error) {
	return as.accountRepo.UpdateAccount(a)
}

func (as *Service) DeleteAccount(id string) error {
	return as.accountRepo.DeleteAccount(id)
}

func (as *Service) GetAccount(id string) (Account, error) {
	return as.accountRepo.GetAccount(id)
}

func (as *Service) AddEvent(e AccountEvent) error {
	_, err := as.GetAccount(e.AccountId)
	if err != nil {
		return err
	}

	_, err = as.locationRepo.GetLocation(e.Location)
	if err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	e.Id = id.String()
	e.Time = time.Now().Unix()
	return as.accountRepo.AddEvent(e)
}
