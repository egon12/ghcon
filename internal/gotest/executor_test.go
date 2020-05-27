package gotest

import (
	"io/ioutil"
	"testing"

	"github.com/egon12/ghr/reviewing"
)

func TestExecutor_Name(t *testing.T) {
	te := NewExecutor()
	if te.Name() != "test-v0.1" {
		t.Error("Wrong name")
	}
}

func TestExecutor_Run(t *testing.T) {
	t.Skip("This test will run test again recursively")

	te := NewExecutor()
	stdout, stderr, _ := te.Run(reviewing.PR{})

	out, _ := ioutil.ReadAll(stdout)
	err, _ := ioutil.ReadAll(stderr)

	t.Log("Out:\n" + string(out))
	t.Log("Err:\n" + string(err))
}
