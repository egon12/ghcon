package main

import (
	"fmt"
	"os"

	"github.com/egon12/ghcon/commit"
	"github.com/egon12/ghcon/path"
	"github.com/egon12/ghcon/review"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const (
	GithubToken = "GITHUB_TOKEN"
)

func main() {
	comment := os.Args[2]

	filePath, line, err := path.ParseFileAndLine(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv(GithubToken)})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	clientV4 := githubv4.NewClient(tc)

	cs := commit.NewSource(clientV4)
	c := cs.GetCurrentCommit()

	r := review.NewProcess(clientV4)
	r.Start(c)

	r.AddComment(filePath, line, comment)

}
