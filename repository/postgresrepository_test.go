package repository

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestSaveTicketPostgres(t *testing.T) {
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_USER", "postgres")
	os.Setenv("POSTGRES_PASS", "password")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_DB", "postgres")
	repo, err := NewRepository("postgres")
	if err != nil {
		t.Fatalf("init map repo: %v\n", err)
	}

	t1 := Ticket{
		Id:       "test",
		Location: "test",
		Start:    time.Now().Unix(),
		Stop:     time.Now().Unix(),
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

func TestGetTicketPostgres(t *testing.T) {
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_USER", "postgres")
	os.Setenv("POSTGRES_PASS", "password")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_DB", "postgres")
	repo, err := NewRepository("postgres")
	if err != nil {
		t.Fatalf("init map repo: %v\n", err)
	}

	now := time.Now().Unix()
	id := fmt.Sprintf("%d", rand.Intn(1000))
	t1 := Ticket{
		Id:       id,
		Location: "test",
		Start:    now,
		Stop:     now,
		Payment:  "test",
	}

	t1, err = repo.SaveTicket(t1)
	if err != nil {
		t.Fatalf("save ticket: %v\n", err)
	}

	t2, err := repo.GetTicket("test", id)
	if err != nil {
		t.Fatalf("save ticket: %v\n", err)
	}

	if !reflect.DeepEqual(t1, t2) {
		t.Fatalf("saved ticket isn't equal: wanted %#v but got %#v\n", t1, t2)
	}
}
