package review

import (
	"context"
	"fmt"

	"github.com/egon12/ghr/diff"
	"github.com/egon12/ghr/githubcommit"
	"github.com/egon12/ghr/path"
	"github.com/shurcooL/githubv4"
)

func NewProcess(clientV4 *githubv4.Client) Process {
	return &process{clientV4: clientV4}
}

type process struct {
	clientV4            *githubv4.Client
	commit              githubcommit.Commit
	pullRequestReviewID string
}

func (r *process) Start(commit githubcommit.Commit) error {
	r.commit = commit
	ok, err := r.continueReview(commit)
	if err != nil {
		return err
	}
	if !ok {
		return r.startReview(commit)
	}
	return nil
}

func (r *process) continueReview(commit githubcommit.Commit) (bool, error) {
	var query struct {
		Repository struct {
			PullRequest struct {
				Reviews struct {
					Nodes []struct {
						ID string
					}
				} `graphql:"reviews(last: 1 states: [PENDING])"`
			} `graphql:"pullRequest(number: $pr_number)"`
		} `graphql:"repository(name: $name, owner: $owner)"`
	}

	owner := commit.GetOwner()
	name := commit.GetRepoName()

	err := r.clientV4.Query(context.TODO(), &query, map[string]interface{}{
		"name":      githubv4.String(name),
		"owner":     githubv4.String(owner),
		"pr_number": githubv4.Int(commit.GetPRNumber()),
	})

	ok := len(query.Repository.PullRequest.Reviews.Nodes) > 0

	if ok {
		r.pullRequestReviewID = query.Repository.PullRequest.Reviews.Nodes[0].ID
	}

	return ok, err
}

func (r *process) startReview(commit githubcommit.Commit) error {
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
	}

	err := r.clientV4.Mutate(context.TODO(), &mutation, inputReview, nil)

	r.pullRequestReviewID = mutation.AddPullRequestReview.PullRequestReview.ID

	return err
}

func (r *process) AddComment(filePath string, line int, comment string) error {
	gitFilePath, err := path.GitPath(filePath)
	if err != nil {
		return fmt.Errorf("Get GitPath error %v", err)
	}

	ghcp, err := diff.NewGithubCommentPosition(r.commit, filePath)
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

	commitOID := githubv4.GitObjectID(r.commit.GetHash())

	input := githubv4.AddPullRequestReviewCommentInput{
		PullRequestReviewID: githubv4.NewID(r.pullRequestReviewID),
		Body:                githubv4.String(comment),
		Path:                githubv4.NewString(githubv4.String(gitFilePath)),
		Position:            githubv4.NewInt(githubv4.Int(position)),
		CommitOID:           githubv4.NewGitObjectID(commitOID),
	}

	err = r.clientV4.Mutate(context.TODO(), &mutation, input, nil)

	return err
}

func (r *process) Finish(finalComment string) error {
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

func (r *process) Cancel() error {
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

func (r *process) Approve(finalComment string) error {
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

func (r *process) RequestChanges(finalComment string) error {
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
