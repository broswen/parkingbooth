package location

import (
	"reflect"
	"testing"
)

func TestSaveLocationService(t *testing.T) {

	repo, err := NewMap()
	if err != nil {
		t.Fatalf("init map repo: %v\n", err)
	}

	service, err := NewService(repo)

	l1 := Location{
		Id:          "test",
		Name:        "test",
		Description: "test",
	}

	l2, err := service.SaveLocation(l1)
	if err != nil {
		t.Fatalf("save location: %v\n", err)
	}

	if !reflect.DeepEqual(l1, l2) {
		t.Fatalf("saved location isn't equal: %#v\n", l2)
	}
}

func TestGetLocationService(t *testing.T) {

	repo, err := NewMap()
	if err != nil {
		t.Fatalf("init map repo: %v\n", err)
	}

	service, err := NewService(repo)

	l1 := Location{
		Id:          "test",
		Name:        "test",
		Description: "test",
	}

	_, err = service.SaveLocation(l1)
	if err != nil {
		t.Fatalf("save location: %v\n", err)
	}

	l2, err := service.GetLocation(l1.Id)
	if err != nil {
		t.Fatalf("save location: %v\n", err)
	}

	if !reflect.DeepEqual(l1, l2) {
		t.Fatalf("saved location isn't equal: %#v\n", l2)
	}
}
