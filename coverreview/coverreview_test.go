package coverreview

import (
	"testing"

	"github.com/egon12/ghr/commit"
	"github.com/egon12/ghr/cover"
)

func TestAddSingleCoverageReview(t *testing.T) {
	c := &coverageReviewer{}
	c.multilineCommenter = &mockMultilineCommenter{}
	nics, _ := cover.GetNotCoverage("cover.out")
	l := c.filterNotInCoverage(nics, []string{"coverreview/coverreview.go"})
	for _, i := range l {
		t.Log(i)
	}
}

type mockMultilineCommenter struct{}

func (m *mockMultilineCommenter) Start(_ commit.Commit) error {
	panic("not implemented") // TODO: Implement
}

func (m *mockMultilineCommenter) AddComment(path string, from int, to int, comment string) error {
	panic("not implemented") // TODO: Implement
}
