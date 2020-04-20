package review

import (
	"context"

	"github.com/egon12/ghr/commit"
	"github.com/google/go-github/v31/github"
)

type multilineCommenter struct {
	commit   commit.Commit
	clientV3 *github.Client
}

func NewMultilineCommenter(clientV3 *github.Client) MultilineCommenter {
	return &multilineCommenter{clientV3: clientV3}
}

func (r *multilineCommenter) Start(commit commit.Commit) error {
	r.commit = commit
	return nil
}

func (r *multilineCommenter) AddComment(path string, fromLineNumber, toLineNumber int, comment string) error {
	side := "RIGHT"

	commitHash := r.commit.GetHash()

	pullRequestComment := github.PullRequestComment{
		Body:      &comment,
		Path:      &path,
		StartLine: &fromLineNumber,
		Line:      &toLineNumber,
		Side:      &side,
		CommitID:  &commitHash,
	}

	_, _, err := r.clientV3.PullRequests.CreateComment(
		context.Background(),
		r.commit.GetOwner(),
		r.commit.GetRepoName(),
		r.commit.GetPRNumber(),
		&pullRequestComment,
	)

	return err
}
