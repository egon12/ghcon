package review

import "github.com/egon12/ghr/githubcommit"

type Process interface {
	Start(githubcommit.Commit) error
	AddComment(path string, lineNumber int, comment string) error
	Cancel() error
	Finish(comment string) error
	Approve(comment string) error
	RequestChanges(comment string) error
}

type MultilineCommenter interface {
	Start(githubcommit.Commit) error
	AddComment(path string, from, to int, comment string) error
}

type CommitSource interface {
	GetCurrentCommit() githubcommit.Commit
}

type ProcessFacade interface {
	Comment(pathAndLineNumber, comment string) error
	Cancel() error
	Finish(string) error
	Approve(string) error
	RequestChanges(string) error
}
