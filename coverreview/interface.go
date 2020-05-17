package coverreview

import "github.com/egon12/ghr/githubcommit"

type ReviewProcess interface {
	Start(githubcommit.Commit) error
	AddComment(path string, line int, comment string) error
	Finish(comment string) error
}
