package location

import (
	"reflect"
	"testing"
)

func TestSaveLocation(t *testing.T) {

	repo, err := NewMap()
	if err != nil {
		t.Fatalf("init map repo: %v\n", err)
	}

	l1 := Location{
		Id:          "test",
		Name:        "test",
		Description: "test",
	}

	l2, err := repo.SaveLocation(l1)
	if err != nil {
		t.Fatalf("save location: %v\n", err)
	}

	if !reflect.DeepEqual(l1, l2) {
		t.Fatalf("saved location isn't equal: %#v\n", l2)
	}
}

func TestGetLocation(t *testing.T) {

	repo, err := NewMap()
	if err != nil {
		t.Fatalf("init map repo: %v\n", err)
	}

	l1 := Location{
		Id:          "test",
		Name:        "test",
		Description: "test",
	}

	_, err = repo.SaveLocation(l1)
	if err != nil {
		t.Fatalf("save location: %v\n", err)
	}

	l2, err := repo.GetLocation(l1.Id)
	if err != nil {
		t.Fatalf("save location: %v\n", err)
	}

	if !reflect.DeepEqual(l1, l2) {
		t.Fatalf("saved location isn't equal: %#v\n", l2)
	}
}
