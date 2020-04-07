package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const (
	GithubToken = "GITHUB_TOKEN"
)

func main() {
	_ = godotenv.Load()

	ctx := context.TODO()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(GithubToken)},
	)

	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	clientV4 := githubv4.NewClient(tc)
	/*
		prs, _, err := client.PullRequests.List(ctx, "egon12", "ghcon", nil)
		if err != nil {

		}
	*/

	ch, err := getCurrentCommitHash()
	if err != nil {
		fmt.Printf("Cannot got commit hash, %s\n", err)
	}

	println(ch)
	prs, _, err := client.PullRequests.ListPullRequestsWithCommit(ctx, "egon12", "ghcon", ch, nil)
	//prs, _, err := client.PullRequests.List(ctx, "egon12", "ghcon", nil)

	var nodeID string
	for _, pr := range prs {
		nodeID = pr.GetNodeID()
		println(pr.GetTitle())
		println(pr.GetNodeID())
		println(pr.GetCommitsURL())
		println(pr.GetMergeCommitSHA())
		println(pr.GetState())

	}
	println(nodeID)

	/*
		review := &github.PullRequestReviewRequest{
			NodeID:   nil,
			CommitID: nil,
			Body:     nil,
			Event:    nil,
			Comments: nil,
		}

		comment := &github.PullRequestComment{
			Body:              nil,
			Path:              nil,
			CommitID:          nil,
			Side:              nil,
			OriginalPosition:  nil,
			StartLine:         nil,
			Line:              nil,
			OriginalLine:      nil,
			OriginalStartLine: nil,
			StartSide:         nil,
			OriginalCommitID:  nil,
		}
	*/

	//client.PullRequests.CreateComment(ctx, "egon12", "ghcon", 1, comment)

	//client.PullRequests.CreateReview(ctx, "egon12", "ghcon", 1, review)

	var review struct {
		AddPullRequestReview struct{} `graphql:"addPullRequestReview(input: $input)"`
	}

	inputReview := githubv4.AddPullRequestReviewInput{
		//	PullRequestID:    nil,
		//	CommitOID:        nil,
		//	Body:             nil,
		//	Event:            nil,
		//	Comments:         nil,
		//	ClientMutationID: nil,
		//}
		//{

		//PullRequestReviewID: nodeID,
		Body:      githubv4.NewString("Trying to add comment to pullRequest"),
		CommitOID: githubv4.NewGitObjectID(githubv4.GitObjectID(ch)),
		//Path:             githubv4.NewString("main.go"),
		//Position:         githubv4.NewInt(12),
		ClientMutationID: githubv4.NewString("someid-forpr-1"),
	}

	err = clientV4.Mutate(ctx, &review, inputReview, nil)
	if err != nil {
		fmt.Printf("Error when mutate: %v", err)
	}

	var m struct {
		AddPullRequestReviewComment struct {
			ClientMutationID string `graphql:"clientMutationId"`
		} `graphql:"addPullRequestReviewComment(input: $input)"`
	}

	inputComment := githubv4.AddPullRequestReviewCommentInput{
		//PullRequestReviewID: nodeID,
		Body:             "Trying to add comment to pullRequest",
		CommitOID:        githubv4.NewGitObjectID(githubv4.GitObjectID(ch)),
		Path:             githubv4.NewString("main.go"),
		Position:         githubv4.NewInt(12),
		ClientMutationID: githubv4.NewString("someid-forpr-1"),
	}
	fmt.Printf("%+v", inputComment)

	err = clientV4.Mutate(ctx, &m, inputComment, nil)
	if err != nil {
		fmt.Printf("Error when mutate: %v", err)
	}

	fmt.Printf("%+v", m)
}

func getCurrentCommitHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	return string(output[:len(output)-1]), err
}
