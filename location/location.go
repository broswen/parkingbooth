package location

type Service struct {
	repo LocationRepository
}

func NewService(repo LocationRepository) (*Service, error) {
	return &Service{
		repo: repo,
	}, nil
}

func (s Service) GetLocation(id string) (Location, error) {
	l, err := s.repo.GetLocation(id)
	if err != nil {
		return Location{}, err
	}
	return l, nil
}

func (s Service) SaveLocation(l Location) (Location, error) {
	l, err := s.repo.SaveLocation(l)
	if err != nil {
		return Location{}, err
	}
	return l, nil
}
