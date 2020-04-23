package coverreview

import (
	"github.com/egon12/ghr/commit"
	"github.com/egon12/ghr/coverage"
	"github.com/egon12/ghr/diff"
	"github.com/egon12/ghr/path"
	"github.com/egon12/ghr/review"
)

type coverageReviewer struct {
	gitDiffProducer    diff.GitDiffProducer
	multilineCommenter review.MultilineCommenter
}

type CoverProfile struct {
	file       string
	percentage float32
	ranges     []coverage.Range
}

func (c *coverageReviewer) Do(commit commit.Commit, coverProfilePath string) error {
	coverProfiles, err := coverage.FromProfile(coverProfilePath)
	if err != nil {
		return err
	}

	listChanges := diff.FromCommit(commit)
	if err != nil {
		return err
	}

	var filteredCoverProfile []CoverProfile
	for _, f := range listChanges.Files() {
		percent := coverProfiles.PercentageFile(f)
		if percent < 1 {
			filteredCoverProfile = append(filteredCoverProfile, CoverProfile{
				file:       f,
				percentage: percent,
				ranges:     coverProfiles.NotInCoverageLines(f),
			})
		}
	}

	return c.AddCoverageReview(commit, filteredCoverProfile)
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
