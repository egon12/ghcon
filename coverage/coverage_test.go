package coverage

import (
	"testing"

	"github.com/egon12/ghr/path"
)

func TestGoCoverageInGit_NotInCoverage(t *testing.T) {
	c, _ := FromProfile("cover.out")
	filepath, _ := path.GetFileWithPackagePath("coverage.go")
	ranges := c.NotInCoverageLines(filepath)

	if len(ranges) == 0 {
		t.Error("It should not zero")
	}
}
