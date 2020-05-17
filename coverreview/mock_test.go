package coverreview

import (
	"github.com/egon12/ghr/commit"
	"github.com/egon12/ghr/coverage"
	"github.com/egon12/ghr/diff"
)

type mockMultilineCommenter struct {
	commit commit.Commit
}

func (m *mockMultilineCommenter) Start(c commit.Commit) error {
	m.commit = c
	return nil
}

func (m *mockMultilineCommenter) AddComment(path string, from int, to int, comment string) error {
	panic("not implemented") // TODO: Implement
}

type mockCoverage struct{}

func (m *mockCoverage) Percentage() float32 {
	panic("not implemented") // TODO: Implement }
}

func (m *mockCoverage) PercentagePackage(packageName string) float32 {
	panic("not implemented") // TODO: Implement
}

func (m *mockCoverage) PercentageFile(filename string) float32 {
	panic("not implemented") // TODO: Implement
}

func (m *mockCoverage) NotInCoverageLines(filename string) []coverage.Range {
	panic("not implemented") // TODO: Implement
}

type mockListChanges struct {
	files []string
}

func (m *mockListChanges) Files() []string {
	return m.files
}

func (m *mockListChanges) RangesInNew(filename string) []diff.Range {
	panic("not implemented") // TODO: Implement
}

func (m *mockListChanges) RangesInOri(filename string) []diff.Range {
	panic("not implemented") // TODO: Implement
}

type mockReviewProcess struct {
}

func (m *mockReviewProcess) Start(_ commit.Commit) error {
	panic("not implemented") // TODO: Implement
}

func (m *mockReviewProcess) AddComment(path string, lineNumber int, comment string) error {
	panic("not implemented") // TODO: Implement
}

func (m *mockReviewProcess) Cancel() error {
	panic("not implemented") // TODO: Implement
}

func (m *mockReviewProcess) Finish(comment string) error {
	panic("not implemented") // TODO: Implement
}

func (m *mockReviewProcess) Approve(comment string) error {
	panic("not implemented") // TODO: Implement
}

func (m *mockReviewProcess) RequestChanges(comment string) error {
	panic("not implemented") // TODO: Implement
}
