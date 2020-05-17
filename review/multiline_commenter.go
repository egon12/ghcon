package review

import (
	"context"
	"fmt"

	"github.com/egon12/ghr/commit"
	"github.com/egon12/ghr/path"
	"github.com/google/go-github/v31/github"
)

type ultilineCommenter struct {
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

func (r *multilineCommenter) AddComment(filePath string, fromLineNumber, toLineNumber int, comment string) error {
	side := "RIGHT"

	commitHash := r.commit.GetHash()

	gitFilePath, err := path.GitPath(filePath)
	if err != nil {
		return fmt.Errorf("Get GitPath error %v", err)
	}

	pullRequestComment := github.PullRequestComment{
		Body:      &comment,
		Path:      &gitFilePath,
		StartLine: &fromLineNumber,
		Line:      &toLineNumber,
		Side:      &side,
		CommitID:  &commitHash,
	}

	_, _, err = r.clientV3.PullRequests.CreateComment(
		context.Background(),
		r.commit.GetOwner(),
		r.commit.GetRepoName(),
		r.commit.GetPRNumber(),
		&pullRequestComment,
	)

	return err
}
