package githubcommit

import (
	"context"
	"os"
	"testing"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/transport/http"
	"gopkg.in/h2non/gock.v1"
)

const repo = "egon12/ghr"

const hash = "cfa9a97b93a300785746b692f6a0de4e7b18aa70"

const searchResponse = `
{"data":{"search":{"nodes":[{"baseRefName":"master","id":"MDExOlB1bGxSZXF1ZXN0Mzg4NzQzODg5","number":1}]}}}`

func TestSource_GetCommitFromRepo_WithMock(t *testing.T) {
	defer gock.Off()
	gock.New("https://api.github.com").
		Post("/graphql").
		Reply(200).
		BodyString(searchResponse)

	client, _, _ := http.NewClient(context.TODO(), option.WithUserAgent("ghr/v0.1.0"))
	clientV4 := githubv4.NewClient(client)
	cs := &Source{clientV4}

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

	if c.GetBaseRefName() != "master" {
		t.Error("BaseRefName is different")
	}

	if c.GetOwner() != "egon12" {
		t.Error("Owner is different")
	}

	if c.GetRepoName() != "ghr" {
		t.Errorf("RepoName is different got %s", c.GetRepoName())
	}
}

func TestSource_GetCommitFromRepo_WithMock_NotPullRequest(t *testing.T) {
	defer gock.Off()
	gock.New("https://api.github.com").
		Post("/graphql").
		Reply(200).
		BodyString(`{"data":{"search":{"nodes":[]}}}`)

	client, _, _ := http.NewClient(context.TODO(), option.WithUserAgent("ghr/v0.1.0"))
	clientV4 := githubv4.NewClient(client)
	cs := &Source{clientV4}

	c := cs.GetCommitFromRepo(repo, hash)

	if c.Error() != nil {
		t.Errorf("Unexpected Error %v", c.Error())
	}

	if c.GetHash() != hash {
		t.Error("hash is different")
	}

	if c.IsPR() == true {
		t.Error("commit should not be a PR")
	}

	if c.GetPRNumber() != 0 {
		t.Error("PR Numbers is different")
	}

	if c.GetPRID() != "" {
		t.Error("PR ID is different")
	}

	if c.GetBaseRefName() != "" {
		t.Error("BaseRefName is different")
	}

	if c.GetOwner() != "egon12" {
		t.Error("Owner is different")
	}

	if c.GetRepoName() != "ghr" {
		t.Errorf("RepoName is different got %s", c.GetRepoName())
	}
}

func Test_getGithubRepoFromOrigin(t *testing.T) {
	origin := "origin\tgit@github.com:egon12/ghr (fetch)\norigin\tgit@github.com:egon12/ghr (push)"

	got, _ := getGithubRepoFromOrigin(origin)
	if got != repo {
		t.Errorf("Want egon12/ghr got %s", got)
	}
}

func Test_getGithubRepoFromOrigin_Failed(t *testing.T) {
	origin := "origin\tgit@gitlab.com:egon12/ghr (fetch)\norigin\tgit@gitlab.com:egon12/ghr (push)"
	want := "Can't find github in 'git remote -v'"

	_, err := getGithubRepoFromOrigin(origin)
	if err.Error() != want {
		t.Errorf("Want %s got %v", want, err)
	}
}

func TestGetCommitFromRepo(t *testing.T) {
	t.Skip("Need GITHUB_TOKEN")

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
	tc := oauth2.NewClient(context.TODO(), ts)

	clientV4 := githubv4.NewClient(tc)

	cs := &Source{clientV4}

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

	if c.GetBaseRefName() != "master" {
		t.Error("BaseRefName is different")
	}
}
