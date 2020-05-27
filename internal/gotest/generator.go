package gotest

import (
	"io"
	"io/ioutil"

	"github.com/egon12/ghr/reviewing"
)

// Generator struct that will convert stdout & stderr
// from gotest into comments
type Generator struct{}

func ProvideNewGenerator() reviewing.Generator {
	return NewGenerator()
}

// NewTestGenerator Create new Test Comment Generator
func NewGenerator() *Generator {
	return &Generator{}
}

func (t *Generator) Generate(stdout, stderr io.Reader, exitCode int) (*reviewing.Review, error) {
	var result reviewing.Review

	b, _ := ioutil.ReadAll(stderr)
	if len(b) > 0 {
		result = reviewing.Review{
			State:   reviewing.Reject,
			Message: "Test Build Error:\n```\n" + string(b) + "\n```\n",
		}
		return &result, nil
	}

	haveError, content, err := readStdOut(stdout)
	if err != nil {
		return &result, nil
	}

	if haveError {
		result = reviewing.Review{
			State:   reviewing.Reject,
			Message: content,
		}
		return &result, nil
	}

	return nil, nil
}
