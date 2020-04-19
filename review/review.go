package review

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/egon12/ghr/diff"
	"github.com/google/go-github/v31/github"
	"github.com/shurcooL/githubv4"
)

func NewProcess(clientV4 *githubv4.Client) Process {
	return &process{clientV4: clientV4}
}

type process struct {
	clientV4            *githubv4.Client
	clientV3            *github.Client
	commit              Commit
	pullRequestReviewID string
}

func (r *process) Start(commit Commit) error {
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

func (r *process) continueReview(commit Commit) (bool, error) {
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

	repoArr := strings.Split(commit.GetRepo(), "/")
	owner := repoArr[0]
	name := repoArr[1]

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

func (r *process) startReview(commit Commit) error {
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

func (r *process) AddComment(path string, line int, comment string) error {
	ghcp, err := diff.NewGithubCommentPosition(r.commit, path)
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
		Path:                githubv4.NewString(githubv4.String(path)),
		Position:            githubv4.NewInt(githubv4.Int(position)),
		CommitOID:           githubv4.NewGitObjectID(commitOID),
	}

	err = r.clientV4.Mutate(context.TODO(), &mutation, input, nil)

	return err
}

func (r *process) AddMultilineComment(path string, fromLineNumber, toLineNumber int, comment string) error {
	var (
		pullRequestReviewID int64
		err                 error
	)

	pullRequestReviewID, err = strconv.ParseInt(r.pullRequestReviewID, 10, 64)
	if err != nil {
		return fmt.Errorf("Cannot get PullRequestReviewID: %v", err)
	}

	side := "right"

	commitHash := r.commit.GetHash()

	ghcp, err := diff.NewGithubCommentPosition(r.commit, path)
	if err != nil {
		return err
	}

	startLine, err := ghcp.Find(fromLineNumber, true)
	if err != nil {
		return err
	}

	line, err := ghcp.Find(toLineNumber, true)
	if err != nil {
		return err
	}

	pullRequestComment := github.PullRequestComment{
		Body:                &comment,
		Path:                &path,
		PullRequestReviewID: &pullRequestReviewID,
		StartLine:           &startLine,
		Line:                &line,
		Side:                &side,
		CommitID:            &commitHash,
	}

	repo := strings.Split(r.commit.GetRepo(), "/")

	_, _, err = r.clientV3.PullRequests.CreateComment(
		context.Background(),
		repo[0], // owner
		repo[1], // name
		r.commit.GetPRNumber(),
		&pullRequestComment,
	)

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
