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
		Message  string
		Comments []Comment
	}

	Source string

	Comment struct {
		Source
		Side    string
		Message string
	}
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
		Finish(State, string) error
		Cancel() error

		// Maybe later Do Start, Comment, Finish
		//Review(Review) error
	}

	Generator interface {
		Generate(stdout, stderr io.Reader, returnCode int) (*Review, error)
	}

	Executor interface {
		Name() string
		Run(pr PR, arguments ...string) (stdout, stderr io.Reader, exitCode int)
	}
)
