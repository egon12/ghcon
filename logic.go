package main

type ReviewProcess interface {
	// StartReview id should be pullRequestNumber, or pullRequest with message
	StartReview(Commit) error
	AddComment(comment, path string, linenumber int) error
	FinishReview() error
	CancelReview() error
}

type ReviewState interface {
	GetOwner() (string, error)

	GetPullRequestID() (string, error)
	GetPullRequestNumber() (int, error)
}

type Commit interface {
	GetRepo() string
	GetHash() string
	IsPR() bool
	GetPRNumber() int
	GetPRID() string
	Error() error
}
