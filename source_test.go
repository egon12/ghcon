package main

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func TestGetCurrentCommit(t *testing.T) {
	t.Skip()
	cs := &commitSource{}
	c := cs.GetCurrentCommit()

	t.Errorf("%#v\n", c)
}

func TestGetGithubRepo(t *testing.T) {
	origin := "origin\tgit@github.com:egon12/ghcon (fetch)\norigin\tgit@github.com:egon12/ghcon (push)"

	got := getGithubRepoFromOrigin(origin)
	if "egon12/ghcon" != got {
		t.Errorf("Want egon12/ghcon got %s", got)
	}
}

func TestGetCommitFromRepo(t *testing.T) {
	repo := "egon12/ghcon"
	hash := "cfa9a97b93a300785746b692f6a0de4e7b18aa70"

	_ = godotenv.Load()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(GithubToken)},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	clientV4 := githubv4.NewClient(tc)

	cs := &commitSource{clientV4}

	c := cs.GetCommitFromRepo(repo, hash)

	if c.Error() != nil {
		t.Errorf("Unexpected Error %v", c.Error())
	}

	if c.GetHash() != hash {
		t.Error("hash is different")
	}

	if c.IsPR() == false {
		t.Error("commit should be a PR")
	}

	if c.GetPRNumber() != 1 {
		t.Error("PR Numbers is different")
	}

	if c.GetPRID() != "MDExOlB1bGxSZXF1ZXN0Mzg4NzQzODg5" {
		t.Error("PR ID is different")
	}
}
