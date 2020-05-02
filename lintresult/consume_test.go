package lintresult

import "testing"

func TestRead(t *testing.T) {
	a, _ := Read("full")
	if len(a) != 10 {
		t.Errorf("Want 10 lint error got %d", len(a))
	}
}
