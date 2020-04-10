package diff

import (
	"fmt"
	"os/exec"
)

type Commit interface {
	GetHash() string
	GetBaseRefName() string
}

type gitDiffProducer struct{}

func (g *gitDiffProducer) Produce(commit Commit, path string) ([]byte, error) {
	branches := fmt.Sprintf("%s..%s", commit.GetHash(), commit.GetBaseRefName())
	cmd := exec.Command("git", "diff", branches, "--", path)
	return cmd.Output()
}
