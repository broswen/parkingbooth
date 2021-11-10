package repository

import (
	"reflect"
	"testing"
	"time"
)

func TestSaveTicket(t *testing.T) {

	repo, err := NewRepository("map")
	if err != nil {
		t.Fatalf("init map repo: %v\n", err)
	}

	t1 := Ticket{
		Id:       "test",
		Location: "test",
		Start:    time.Now(),
		Stop:     time.Now(),
		Payment:  "test",
	}

	t2, err := repo.SaveTicket(t1)
	if err != nil {
		t.Fatalf("save ticket: %v\n", err)
	}

	if !reflect.DeepEqual(t1, t2) {
		t.Fatalf("saved ticket isn't equal: %#v\n", t2)
	}
}

func TestGetTicket(t *testing.T) {

	repo, err := NewRepository("map")
	if err != nil {
		t.Fatalf("init map repo: %v\n", err)
	}

	t1 := Ticket{
		Id:       "test",
		Location: "test",
		Start:    time.Now(),
		Stop:     time.Now(),
		Payment:  "test",
	}

	_, err = repo.SaveTicket(t1)
	if err != nil {
		t.Fatalf("save ticket: %v\n", err)
	}

	t2, err := repo.GetTicket("test", "test")
	if err != nil {
		t.Fatalf("save ticket: %v\n", err)
	}

	if !reflect.DeepEqual(t1, t2) {
		t.Fatalf("saved ticket isn't equal: %#v\n", t2)
	}
}
