package diff

import (
	"fmt"
	"os/exec"
	"strings"
)

type Commit interface {
	GetHash() string
	GetBaseRefName() string
}

type gitDiffProducer struct{}

func (g *gitDiffProducer) Produce(commit Commit, path string) ([]byte, error) {
	branches := g.getBranches(commit)
	cmd := exec.Command("git", "diff", branches, "--", path)
	return cmd.Output()
}

func (g *gitDiffProducer) ListFiles(commit Commit) ([]string, error) {
	branches := g.getBranches(commit)
	cmd := exec.Command("git", "diff", "--stat", branches)
	b, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	result := strings.Split(string(b), "\n")

	for i, line := range result {
		cell := strings.Split(line, "|")
		result[i] = strings.TrimSpace(cell[0])
	}
	return result, err
}

func (g *gitDiffProducer) getBranches(commit Commit) string {
	return fmt.Sprintf("%s...%s", commit.GetBaseRefName(), commit.GetHash())
}
