package lintresult

import "testing"

func TestRead(t *testing.T) {
	a, _ := Read("full")
	t.Errorf("%+v", a)
}
