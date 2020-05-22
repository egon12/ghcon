package githubcommit

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/shurcooL/githubv4"
)

// Source are commit source, or we can say as commit repository.
type Source struct {
	githubClientV4 *githubv4.Client
}

// NewSource create new commit source.
func NewSource(client *githubv4.Client) *Source {
	return &Source{client}
}

// GetCurrentCommit are get commit data from directory and what you checkout.
func (c *Source) GetCurrentCommit() Commit {
	cmd := exec.Command("git", "rev-parse", "HEAD")

	output, err := cmd.Output()
	if err != nil {
		return &commit{err: err}
	}

	hash := string(output[:len(output)-1])

	return c.GetCommit(hash)
}

// GetCommit are function to get commit data from your directory.
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

// GetCommitFromRepo are function to get commit data with repo and hash.
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

	return getGithubRepoFromOrigin(origin)
}

func getGithubRepoFromOrigin(origin string) (string, error) {
	r := regexp.MustCompile(`github.com[:/]([a-zA-Z0-9/_]*)`)

	m := r.FindStringSubmatch(origin)
	if m == nil {
		return "", errors.New("Can't find github in 'git remote -v'")
	}

	return m[1], nil
}
