package main

import (
	"fmt"
	"os"

	"github.com/egon12/ghr/commit"
	"github.com/egon12/ghr/path"
	"github.com/egon12/ghr/review"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const (
	GithubToken = "GITHUB_TOKEN"
)

func main() {
	if os.Args[1] == "cancel" {
		cancel()
		return
	}

	if os.Args[1] == "finish" {
		finish()
		return
	}

	addComment()
}

func start() {
	cs := getCommitSource()
	r := getReviewProcess()

	c := cs.GetCurrentCommit()

	err := r.Start(c)
	if err != nil {
		fmt.Fprint(os.Stdout, err)
		os.Exit(1)
	}
}

func finish() {
	start()
	comment := os.Args[2]
	r := getReviewProcess()
	r.Finish(comment)
}

func cancel() {
	start()
	r := getReviewProcess()
	r.Cancel()
}

func addComment() {
	start()
	comment := os.Args[2]

	filePath, line, err := path.ParseFileAndLine(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	r := getReviewProcess()
	err = r.AddComment(filePath, line, comment)
	if err != nil {
		fmt.Fprint(os.Stdout, err)
	}
}

var clientV4 *githubv4.Client

func getClient() *githubv4.Client {
	if clientV4 == nil {
		githubToken := os.Getenv(GithubToken)
		if githubToken == "" {
			fmt.Fprintf(os.Stdout, "Please set GITHUB_TOKEN environment variable")
			os.Exit(1)
		}
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
		tc := oauth2.NewClient(oauth2.NoContext, ts)
		clientV4 = githubv4.NewClient(tc)
	}
	return clientV4
}

func getCommitSource() *commit.Source {
	return commit.NewSource(getClient())
}

var reviewProcess review.Process

func getReviewProcess() review.Process {
	if reviewProcess == nil {
		reviewProcess = review.NewProcess(getClient())
	}
	return reviewProcess
}
