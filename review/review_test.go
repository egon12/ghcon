package review

import (
	"os"
	"testing"

	"github.com/egon12/ghr/commit"
	"github.com/joho/godotenv"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const (
	GithubToken = "GITHUB_TOKEN"
)

func getClient() *githubv4.Client {
	_ = godotenv.Load("../.env")
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv(GithubToken)})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return githubv4.NewClient(tc)
}

func TestReviewProcess_StartAddCommentFinish(t *testing.T) {
	t.Skip()
	hash := "cfa9a97b93a300785746b692f6a0de4e7b18aa70"
	clientV4 := getClient()
	cs := commit.NewSource(clientV4)
	r := &process{clientV4: clientV4}

	err := r.Start(cs.GetCommit(hash))
	if err != nil {
		t.Error(err)
	}

	err = r.AddComment("main.go", 14, "Try to comment at GithubToken main.go:14")
	if err != nil {
		t.Error(err)
	}

	err = r.Finish("And this is final comment")
	if err != nil {
		t.Error(err)
	}
}

func TestReviewProcess_StartReview(t *testing.T) {
	t.Skip()
	hash := "cfa9a97b93a300785746b692f6a0de4e7b18aa70"
	clientV4 := getClient()
	cs := commit.NewSource(clientV4)
	r := &process{clientV4: clientV4}

	err := r.Start(
		cs.GetCommit(hash),
	)

	if err != nil {
		t.Error(err)
	}
}

func TestReviewProcess_AddComment(t *testing.T) {
	t.Skip()
	r := &process{
		clientV4:            getClient(),
		pullRequestReviewID: "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3Mzg5NDE4MTA5",
	}

	err := r.AddComment("main.go", 14, "Try to comment at GithubToken main.go:14")
	if err != nil {
		t.Error(err)
	}
}

func TestReviewProcess_AddMultilineComment(t *testing.T) {
	t.Skip()
	r := &process{
		clientV4:            getClient(),
		pullRequestReviewID: "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3Mzg5NDE4MTA5",
	}

	err := r.AddComment("main.go", 14, "Try to comment at GithubToken main.go:14")
	if err != nil {
		t.Error(err)
	}
}

func TestReviewProcess_CancelReview(t *testing.T) {
	t.Skip()
	r := &process{
		clientV4:            getClient(),
		pullRequestReviewID: "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3Mzg5NDE4MTA5",
	}

	err := r.Cancel()
	if err != nil {
		t.Error(err)
	}
}
