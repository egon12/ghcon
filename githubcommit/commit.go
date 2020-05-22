package githubcommit

import (
	"regexp"
	"strings"
)

// Commit is commit data
type Commit interface {
	GetRepo() string
	GetOwner() string
	GetRepoName() string
	GetHash() string
	IsPR() bool
	GetPRNumber() int
	GetPRID() string
	GetBaseRefName() string
	Error() error
}

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
	r := regexp.MustCompile(`([a-zA-Z0-9_-]+)/([a-zA-Z0-9_-]+)`)

	m := r.FindStringSubmatch(c.repo)
	if m == nil {
		return ""
	}

	return m[2]
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
