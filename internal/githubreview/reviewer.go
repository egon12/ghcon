package githubreview

import (
	"context"
	"fmt"

	"github.com/egon12/ghr/diff"
	"github.com/egon12/ghr/path"
	"github.com/egon12/ghr/reviewing"
	"github.com/shurcooL/githubv4"
)

type Reviewer struct {
	clientV4 *githubv4.Client
	commit   reviewing.Commit
	pr       reviewing.PR
	reviewID string
}

// ProvideReviewer implement the revieweing.Reviewer
func ProvideReviewer(c *githubv4.Client) reviewing.Reviewer {
	return NewReviewer(c)
}

func NewReviewer(clientV4 *githubv4.Client) *Reviewer {
	return &Reviewer{clientV4: clientV4}
}

func (r *Reviewer) Start(pr reviewing.PR, commit reviewing.Commit) error {
	r.pr = pr
	r.commit = commit

	var mutation struct {
		AddPullRequestReview struct {
			PullRequestReview struct {
				ID string
			}
		} `graphql:"addPullRequestReview(input: $input)"`
	}

	commitOID := githubv4.GitObjectID(commit.Hash)

	inputReview := githubv4.AddPullRequestReviewInput{
		PullRequestID: githubv4.NewID(pr.ID),
		CommitOID:     githubv4.NewGitObjectID(commitOID),
	}

	err := r.clientV4.Mutate(context.TODO(), &mutation, inputReview, nil)

	r.reviewID = mutation.AddPullRequestReview.PullRequestReview.ID

	return err
}

func (r *Reviewer) Comment(comment reviewing.Comment) error {
	gitFilePath, err := path.GitPath(comment.Source)
	if err != nil {
		return fmt.Errorf("Get GitPath error %v", err)
	}

	ghcp, err := diff.NewGithubCommentPosition(r.commit, comment.Source)
	if err != nil {
		return err
	}

	position, err := ghcp.Find(line, true)
	if err != nil {
		return err
	}

	var mutation struct {
		AddPullRequestReviewComment struct {
			ClientMutationID string
		} `graphql:"addPullRequestReviewComment(input: $input)"`
	}

	commitOID := githubv4.GitObjectID(r.commit.Hash)

	input := githubv4.AddPullRequestReviewCommentInput{
		PullRequestReviewID: githubv4.NewID(r.reviewID),
		Body:                githubv4.String(comment.Message),
		Path:                githubv4.NewString(githubv4.String(gitFilePath)),
		Position:            githubv4.NewInt(githubv4.Int(position)),
		CommitOID:           githubv4.NewGitObjectID(commitOID),
	}

	err = r.clientV4.Mutate(context.TODO(), &mutation, input, nil)

	return err
}

func (r *Reviewer) Finish(comment reviewing.Comment) error {
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

	err := r.clientV4.Mutate(context.TODO(), &mutation, input, nil)

	return err
}

func (r *Reviewer) Cancel() error {
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

	err := r.clientV4.Mutate(context.TODO(), &mutation, input, nil)
	if err == nil && mutation.DeletePullRequestReview.ClientMutationID != mutationID {
		return fmt.Errorf("Got different mutationID want %s got %s", mutationID, mutation.DeletePullRequestReview.ClientMutationID)
	}

	return err
}

/*


func (r *Reviewer) Approve(finalComment string) error {
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
		Event:               githubv4.PullRequestReviewEventApprove,
		Body:                body,
	}

	err := r.clientV4.Mutate(context.TODO(), &mutation, input, nil)

	return err
}

func (r *Reviewer) RequestChanges(finalComment string) error {
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
		Event:               githubv4.PullRequestReviewEventRequestChanges,
		Body:                body,
	}

	err := r.clientV4.Mutate(context.TODO(), &mutation, input, nil)

	return err
}
*/
