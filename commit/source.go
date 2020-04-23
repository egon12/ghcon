package commit

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/shurcooL/githubv4"
)

type commit struct {
	repo         string
	hash         string
	err          error
	isPR         bool
	prNumbers    []int
	prID         []string
	baseRefNames []string
}

func (c *commit) GetRepo() string {
	return c.repo
}

func (c *commit) GetOwner() string {
	ownerAndRepo := strings.Split(c.repo, "/")
	return ownerAndRepo[0]
}

func (c *commit) GetRepoName() string {
	ownerAndRepo := strings.Split(c.repo, "/")
	if len(ownerAndRepo) < 2 {
		return ""
	}
	return ownerAndRepo[1]
}

func (c *commit) GetHash() string {
	return c.hash
}

func (c *commit) Error() error {
	return c.err
}

func (c *commit) IsPR() bool {
	return c.isPR
}

func (c *commit) GetPRNumber() int {
	if len(c.prNumbers) > 0 {
		return c.prNumbers[0]
	}
	return 0
}

func (c *commit) GetPRID() string {
	if len(c.prID) > 0 {
		return c.prID[0]
	}
	return ""
}

func (c *commit) GetBaseRefName() string {
	if len(c.baseRefNames) > 0 {
		return c.baseRefNames[0]
	}
	return ""
}

type Source struct {
	githubClientV4 *githubv4.Client
}

func NewSource(client *githubv4.Client) *Source {
	return &Source{client}
}

func (c *Source) GetCurrentCommit() Commit {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return &commit{err: err}
	}

	hash := string(output[:len(output)-1])
	return c.GetCommit(hash)
}

func (c *Source) GetCommit(hash string) Commit {
	repo, err := getGithubRepo()
	if err != nil {
		return &commit{
			err: err,
		}
	}

	if repo == "" {
		return &commit{
			err: fmt.Errorf("this git repo not have github origin"),
		}
	}

	return c.GetCommitFromRepo(repo, hash)
}

func (c *Source) GetCommitFromRepo(repo, hash string) Commit {
	result := commit{repo: repo, hash: hash}

	var query struct {
		Search struct {
			Nodes []struct {
				PullRequest struct {
					BaseRefName string
					ID          string
					Number      int
				} `graphql:"... on PullRequest"`
			}
		} `graphql:"search(type: ISSUE, query: $query, last:20)"`
	}

	searchQuery := fmt.Sprintf("repo:%s type:pr %s", repo, hash)
	err := c.githubClientV4.Query(context.TODO(), &query, map[string]interface{}{
		"query": githubv4.String(searchQuery),
	})

	if err != nil {
		result.err = err
		return &result
	}

	prLength := len(query.Search.Nodes)

	result.isPR = prLength > 0
	result.prNumbers = make([]int, prLength)
	result.prID = make([]string, prLength)
	result.baseRefNames = make([]string, prLength)

	for i, pr := range query.Search.Nodes {
		result.prNumbers[i] = pr.PullRequest.Number
		result.prID[i] = pr.PullRequest.ID
		result.baseRefNames[i] = pr.PullRequest.BaseRefName
	}

	return &result
}

func getGithubRepo() (string, error) {
	cmd := exec.Command("git", "remote", "-v")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	origin := string(output[:len(output)-1])
	return getGithubRepoFromOrigin(origin), nil
}

func getGithubRepoFromOrigin(origin string) string {
	r := regexp.MustCompile(`github.com[:/]([a-zA-Z0-9/_]*)`)

	m := r.FindAllStringSubmatch(origin, 1)
	if len(m) < 1 {
		return ""
	}

	if len(m[0]) < 2 {
		return ""
	}

	return m[0][1]
}
