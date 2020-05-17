package coverage

import (
	"testing"
)

func TestGoCoverageInGit_NotInCoverage(t *testing.T) {
	c, _ := FromProfile("cover.out")
	ranges := c.NotInCoverageLines("coverage.go")

	if len(ranges) == 0 {
		t.Error("It should not zero")
	}
}
