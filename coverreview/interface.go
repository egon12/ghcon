package coverreview

import "github.com/egon12/ghr/commit"

type ReviewProcess interface {
	Start(commit.Commit) error
	AddComment(path string, line int, comment string) error
	Finish(comment string) error
}
