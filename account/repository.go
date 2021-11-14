package account

import "fmt"

type Account struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AccountEventType string

const (
	InEvent  AccountEventType = "IN"
	OutEvent AccountEventType = "OUT"
)

type AccountEvent struct {
	Id        string           `json:"id"`
	AccountId string           `json:"accountId"`
	Type      AccountEventType `json:"type"`
	Location  string           `json:"location"`
	Time      int64            `json:"time"`
}

type AccountRepository interface {
	GetAccount(id string) (Account, error)
	AddEvent(e AccountEvent) error
	CreateAccount(a Account) (Account, error)
	UpdateAccount(a Account) (Account, error)
	DeleteAccount(a Account) error
}

type MapRepository struct {
	m map[string]Account
	e map[string][]AccountEvent
}

func NewMap() (AccountRepository, error) {
	return MapRepository{
		m: make(map[string]Account, 0),
		e: make(map[string][]AccountEvent, 0),
	}, nil
}

func (mr MapRepository) CreateAccount(a Account) (Account, error) {
	_, ok := mr.m[a.Id]
	if ok {
		return Account{}, fmt.Errorf("account with id already exists")
	}
	mr.m[a.Id] = a
	mr.e[a.Id] = make([]AccountEvent, 0)
	return a, nil
}

func (mr MapRepository) GetAccount(id string) (Account, error) {
	a, ok := mr.m[id]
	if !ok {
		return Account{}, fmt.Errorf("account doesn't exist")
	}
	return a, nil
}

func (mr MapRepository) UpdateAccount(a Account) (Account, error) {
	_, ok := mr.m[a.Id]
	if !ok {
		return Account{}, fmt.Errorf("account doesn't exist")
	}
	mr.m[a.Id] = a
	return a, nil
}

func (mr MapRepository) DeleteAccount(a Account) error {
	delete(mr.m, a.Id)
	delete(mr.e, a.Id)
	return nil
}

func (mr MapRepository) AddEvent(e AccountEvent) error {
	events, ok := mr.e[e.AccountId]
	if !ok {
		return fmt.Errorf("account doesn't exist")
	}
	mr.e[e.AccountId] = append(events, e)
	return nil
}
