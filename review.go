package main

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

type reviewProcess struct {
	githubClientV4      *githubv4.Client
	pullRequestReviewID string
}

func (r *reviewProcess) StartReview(commit Commit) error {
	var mutation struct {
		AddPullRequestReview struct {
			PullRequestReview struct {
				ID string
			}
		} `graphql:"addPullRequestReview(input: $input)"`
	}

	commitOID := githubv4.GitObjectID(commit.GetHash())

	inputReview := githubv4.AddPullRequestReviewInput{
		PullRequestID: githubv4.NewID(commit.GetPRID()),
		CommitOID:     githubv4.NewGitObjectID(commitOID),
		Body:          githubv4.NewString("Starting Review"),
	}

	err := r.githubClientV4.Mutate(context.TODO(), &mutation, inputReview, nil)

	r.pullRequestReviewID = mutation.AddPullRequestReview.PullRequestReview.ID

	return err
}

func (r *reviewProcess) AddComment(comment, path string, line int) error {
	var mutation struct {
		AddPullRequestReviewComment struct {
			ClientMutationID string
		} `graphql:"addPullRequestReviewComment(input: $input)"`
	}

	input := githubv4.AddPullRequestReviewCommentInput{
		PullRequestReviewID: githubv4.NewID(r.pullRequestReviewID),
		Body:                githubv4.String(comment),
		Path:                githubv4.NewString(githubv4.String(path)),
		Position:            githubv4.NewInt(githubv4.Int(line)),
	}

	err := r.githubClientV4.Mutate(context.TODO(), &mutation, input, nil)

	return err
}

func (r *reviewProcess) FinishReview(finalComment string) error {
	var mutation struct {
		SubmitPullRequestReview struct {
			ClientMutationID string
		} `graphql:"submitPullRequestReview(input: $input)"`
	}

	var body *githubv4.String
	if finalComment != "" {
		body = githubv4.NewString(githubv4.String(finalComment))
	}

	input := githubv4.SubmitPullRequestReviewInput{
		PullRequestReviewID: githubv4.NewID(r.pullRequestReviewID),
		Event:               githubv4.PullRequestReviewEventComment,
		Body:                body,
	}

	err := r.githubClientV4.Mutate(context.TODO(), &mutation, input, nil)

	return err
}

func (r *reviewProcess) CancelReview() error {
	mutationID := "cancel-" + r.pullRequestReviewID

	var mutation struct {
		DeletePullRequestReview struct {
			ClientMutationID string
		} `graphql:"deletePullRequestReview(input: $input)"`
	}

	input := githubv4.DeletePullRequestReviewInput{
		PullRequestReviewID: githubv4.NewID(r.pullRequestReviewID),
		ClientMutationID:    githubv4.NewString(githubv4.String(mutationID)),
	}

	err := r.githubClientV4.Mutate(context.TODO(), &mutation, input, nil)
	if err == nil && mutation.DeletePullRequestReview.ClientMutationID != mutationID {
		return fmt.Errorf("Got different mutationID want %s got %s", mutationID, mutation.DeletePullRequestReview.ClientMutationID)
	}

	return err

}
