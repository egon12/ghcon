package review

import "github.com/egon12/ghr/commit"

type Process interface {
	// StartReview id should be pullRequestNumber, or pullRequest with message
	Start(commit.Commit) error
	AddComment(path string, lineNumber int, comment string) error
	AddMultilineComment(path string, fromLineNumber, toLineNumber int, comment string) error
	Cancel() error
	Finish(comment string) error
	Approve(comment string) error
	RequestChanges(comment string) error
}

type CommitSource interface {
	GetCurrentCommit() commit.Commit
}

type ProcessFacade interface {
	Comment(path string, lineNumber int, comment string) error
	MultilineComment(path string, fromLineNumber, toLineNumber int, comment string) error
	Cancel() error
	Finish(string) error
	Approve(string) error
	RequestChanges(string) error
}
