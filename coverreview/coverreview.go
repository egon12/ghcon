package coverreview

import (
	"github.com/egon12/ghr/commit"
	"github.com/egon12/ghr/coverage"
	"github.com/egon12/ghr/diff"
	"github.com/egon12/ghr/path"
	"github.com/egon12/ghr/review"
)

type coverageReviewer struct {
	multilineCommenter review.MultilineCommenter
}

type CoverProfile struct {
	file       string
	percentage float32
	ranges     []coverage.Range
}

func (c *coverageReviewer) Do(com commit.Commit, coverProfilePath string) error {
	coverProfiles, err := coverage.FromProfile(coverProfilePath)
	if err != nil {
		return err
	}

	listChanges := diff.FromCommit(com)
	if err != nil {
		return err
	}

	return c.DoReview(com, coverProfiles, listChanges)
}

func (c *coverageReviewer) DoReview(com commit.Commit, cv coverage.GoCoverageInGit, l diff.ListChanges) error {
	var filteredCoverProfile []CoverProfile
	for _, f := range l.Files() {
		percent := cv.PercentageFile(f)
		if percent < 1 {
			filteredCoverProfile = append(filteredCoverProfile, CoverProfile{
				file:       f,
				percentage: percent,
				ranges:     cv.NotInCoverageLines(f),
			})
		}
	}

	return c.AddCoverageReview(com, filteredCoverProfile)
}

func (c *coverageReviewer) AddCoverageReview(commit commit.Commit, coverProfile []CoverProfile) error {
	c.multilineCommenter.Start(commit)
	for _, cp := range coverProfile {
		err := c.AddSingleFileCoverageReview(cp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *coverageReviewer) AddSingleFileCoverageReview(cp CoverProfile) error {
	for _, r := range cp.ranges {
		err := c.AddSingleCoverageReview(cp.file, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *coverageReviewer) AddSingleCoverageReview(filePath string, lineRange coverage.Range) error {
	gitPath, _ := path.GitPath(filePath)
	return c.multilineCommenter.AddComment(gitPath, lineRange.From, lineRange.To, "Not In Coverage")
}
