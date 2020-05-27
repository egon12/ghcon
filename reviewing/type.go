package reviewing

import "io"

type (
	Repo struct {
		Owner string
		Name  string
	}

	Commit struct {
		Hash string
		Repo
	}

	PR struct {
		ID          string
		Number      int
		BaseRefName string
		Reviews     []Review
		Repo
	}

	State int

	Review struct {
		State
		Comments []Comment
	}

	CommentType int

	Comment struct {
		CommentType
		State
		Source  string
		Side    string
		Message string
	}
)

const (
	CommentWithSource CommentType = iota + 1
	CommentWithState
)

const (
	Neutral State = iota + 1
	Approve
	Reject
	Dismiss
)

type (
	GithubSource interface {
		GetCurrentCommit() (Commit, error)
		GetPRFromCommit(Commit) (PR, error)
		GetCurrentPR() (PR, error)
	}

	Reviewer interface {
		Start(PR, Commit) error
		Comment(Comment) error
		Finish(Comment) error
		Cancel() error
	}

	Generator interface {
		Generate(stdout, stderr io.Reader, returnCode int) ([]Comment, error)
	}

	Executor interface {
		Name() string
		Run(pr PR, arguments ...string) (stdout, stderr io.Reader, exitCode int)
	}
)
