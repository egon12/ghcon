package coverreview

import (
	"testing"

	"github.com/egon12/ghr/commit"
)

func TestAddSingleCoverageReview(t *testing.T) {
}

type mockMultilineCommenter struct{}

func (m *mockMultilineCommenter) Start(_ commit.Commit) error {
	panic("not implemented") // TODO: Implement
}

func (m *mockMultilineCommenter) AddComment(path string, from int, to int, comment string) error {
	panic("not implemented") // TODO: Implement
}
