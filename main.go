package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
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
	prs, _, err := client.PullRequests.List(ctx, "egon12", "ghcon", nil)
	if err != nil {

	}

	for _, pr := range prs {
		println(pr.GetState())
	}

	println(client)
	fmt.Println("vim-go")
}
