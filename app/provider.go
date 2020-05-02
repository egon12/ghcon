package app

import (
	"net/http"

	"github.com/egon12/ghr/commit"
	"github.com/egon12/ghr/coverreview"
	"github.com/egon12/ghr/review"
	"github.com/google/go-github/v31/github"
	"github.com/google/wire"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var set = wire.NewSet(
	ProvideGithubToken,
	ProvideOauthTokenSource,
	ProvideOauthClient,
	githubv4.NewClient,
	github.NewClient,
	review.NewProcess,
	review.NewMultilineCommenter,
	review.NewProcessFacade,
	commit.NewSource,
	coverreview.NewCoverageReviewer,
	wire.Bind(new(review.CommitSource), new(*commit.Source)),
	wire.Struct(new(App), "ReviewProcess", "CoverageReviewer"),
)

type GithubToken string

func ProvideOauthTokenSource(githubToken GithubToken) oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{AccessToken: string(githubToken)})
}

func ProvideOauthClient(ts oauth2.TokenSource) *http.Client {
	return oauth2.NewClient(oauth2.NoContext, ts)
}
