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

func (t *Generator) Generate(stdout, stderr io.Reader, exitCode int) ([]reviewing.Comment, error) {
	var result []reviewing.Comment

	b, _ := ioutil.ReadAll(stderr)
	if len(b) > 0 {
		result = append(result, reviewing.Comment{
			CommentType: reviewing.CommentWithState,
			State:       reviewing.Reject,
			Message:     "Test Build Error:\n```\n" + string(b) + "\n```\n",
		})
		return result, nil
	}

	haveError, content, err := readStdOut(stdout)
	if err != nil {
		return result, nil
	}

	if haveError {
		result = append(result, reviewing.Comment{
			CommentType: reviewing.CommentWithState,
			State:       reviewing.Reject,
			Message:     content,
		})
		return result, nil
	}

	return result, nil
}
