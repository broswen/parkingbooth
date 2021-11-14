package account

type Service struct {
	accountRepo AccountRepository
}

func NewService(accountRepo AccountRepository) (*Service, error) {
	return &Service{
		accountRepo: accountRepo,
	}, nil
}

func (as *Service) CreateAccount(a Account) (Account, error) {
	return as.accountRepo.CreateAccount(a)
}

func (as *Service) UpdateAccount(a Account) (Account, error) {
	return as.accountRepo.UpdateAccount(a)
}

func (as *Service) DeleteAccount(a Account) error {
	return as.accountRepo.DeleteAccount(a)
}

func (as *Service) AddEvent(e AccountEvent) error {
	return as.accountRepo.AddEvent(e)
}
