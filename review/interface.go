package review

import "github.com/egon12/ghr/commit"

type Process interface {
	Start(commit.Commit) error
	AddComment(path string, lineNumber int, comment string) error
	Cancel() error
	Finish(comment string) error
	Approve(comment string) error
	RequestChanges(comment string) error
}

type MultilineCommenter interface {
	Start(commit.Commit) error
	AddComment(path string, from, to int, comment string) error
}

type CommitSource interface {
	GetCurrentCommit() commit.Commit
}

type ProcessFacade interface {
	Comment(pathAndLineNumber, comment string) error
	Cancel() error
	Finish(string) error
	Approve(string) error
	RequestChanges(string) error
}
