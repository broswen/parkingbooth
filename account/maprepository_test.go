package account

import (
	"reflect"
	"testing"
)

func TestCreateAccount(t *testing.T) {

	repo, err := NewMap()
	if err != nil {
		t.Fatalf("init map repo: %v\n", err)
	}

	a1 := Account{
		Id:    "test",
		Name:  "test account",
		Email: "contact email",
	}

	_, err = repo.CreateAccount(a1)
	if err != nil {
		t.Fatalf("save account: %v\n", err)
	}

	a2, err := repo.GetAccount(a1.Id)
	if err != nil {
		t.Fatalf("get account: %v\n", err)
	}

	if !reflect.DeepEqual(a1, a2) {
		t.Fatalf("saved account isn't equal: %#v\n", a2)
	}

	_, err = repo.CreateAccount(a1)
	if err == nil {
		t.Fatalf("create with same id should throw error: %#v\n", a1)
	}

	name := "updated name"
	a2.Name = name
	a3, err := repo.UpdateAccount(a2)
	if err != nil {
		t.Fatalf("update account: %#v\n", err)
	}

	if a3.Name != name {
		t.Fatalf("account name not updated: %#v\n", a3.Name)

	}

	err = repo.DeleteAccount(a3)
	if err != nil {
		t.Fatalf("delete account: %#v\n", err)
	}

	_, err = repo.GetAccount(a3.Id)
	if err == nil {
		t.Fatalf("deleted account should throw error: %#v\n", err)

	}
}
