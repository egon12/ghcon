package review

type Commit interface {
	GetRepo() string
	GetHash() string
	IsPR() bool
	GetPRNumber() int
	GetPRID() string
	GetBaseRefName() string
	Error() error
}

type Process interface {
	// StartReview id should be pullRequestNumber, or pullRequest with message
	Start(Commit) error
	AddComment(path string, lineNumber int, comment string) error
	AddMultilineComment(path string, fromLineNumber, toLineNumber int, comment string) error
	Finish(lastComment string) error
	Cancel() error
}
