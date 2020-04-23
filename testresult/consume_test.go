package testresult

import "testing"

func TestReadError(t *testing.T) {
	a, _ := Read("out", "err")
	if len(a) != 532 {
		t.Errorf("len(a) = %d", len(a))
	}
}

func TestReadOut(t *testing.T) {
	a, _ := Read("out", "empty")
	if len(a) != 244 {
		t.Errorf("len(a) = %d", len(a))
	}
}

func TestReadSucces(t *testing.T) {
	a, _ := Read("empty", "empty")
	if len(a) != 0 {
		t.Errorf("len(a) = %d", len(a))
	}
}
