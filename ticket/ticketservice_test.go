package ticket

import (
	"reflect"
	"testing"

	"github.com/broswen/parkingbooth/location"
)

func TestGenerateTicket(t *testing.T) {
	repo, err := NewMap()
	if err != nil {
		t.Fatalf("init map repo: %v\n", err)
	}

	locationRepo, err := location.NewMap()
	if err != nil {
		t.Fatalf("init location map repo: %v\n", err)
	}

	service, err := NewService(repo, locationRepo)
	if err != nil {
		t.Fatalf("init ticket service: %v\n", err)
	}

	t1, err := service.GenerateTicket("test")
	if err != nil {
		t.Fatalf("generate ticket: %v\n", err)
	}

	t2, err := service.GetTicket(t1.Location, t1.Id)
	if err != nil {
		t.Fatalf("get ticket: %v\n", err)
	}

	if !reflect.DeepEqual(t1, t2) {
		t.Fatalf("tickets don't match")
	}

}

func TestCompleteTicket(t *testing.T) {
	repo, err := NewMap()
	if err != nil {
		t.Fatalf("init map repo: %v\n", err)
	}

	locationRepo, err := location.NewMap()
	if err != nil {
		t.Fatalf("init location map repo: %v\n", err)
	}

	service, err := NewService(repo, locationRepo)
	if err != nil {
		t.Fatalf("init ticket service: %v\n", err)
	}

	t1, err := service.GenerateTicket("test")
	if err != nil {
		t.Fatalf("generate ticket: %v\n", err)
	}

	t1, err = service.CompleteTicket(t1.Location, t1.Id)
	if err != nil {
		t.Fatalf("complete ticket: %v\n", err)
	}
}

func TestPayTicket(t *testing.T) {
	repo, err := NewMap()
	if err != nil {
		t.Fatalf("init map repo: %v\n", err)
	}

	locationRepo, err := location.NewMap()
	if err != nil {
		t.Fatalf("init location map repo: %v\n", err)
	}

	service, err := NewService(repo, locationRepo)
	if err != nil {
		t.Fatalf("init ticket service: %v\n", err)
	}

	t1, err := service.GenerateTicket("test")
	if err != nil {
		t.Fatalf("generate ticket: %v\n", err)
	}

	t1, err = service.CompleteTicket(t1.Location, t1.Id)
	if err != nil {
		t.Fatalf("complete ticket: %v\n", err)
	}

	t1, err = service.PayTicket(t1.Location, t1.Id, "test")
	if err != nil {
		t.Fatalf("pay ticket: %v\n", err)
	}
}
