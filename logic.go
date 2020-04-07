package main

type ReviewProcess interface {
	// StartReview id should be pullRequestNumber, or pullRequest with message
	StartReview(hash string) error
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
	GetHash() string
	IsPR() bool
	GetPRNumber() int
	Error() error
}
