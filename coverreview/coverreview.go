package coverreview

// pacakge coverreview
//
// This package will add review about coverage
// While there are too many options for now..
// What are the options?:
//
// # Comment all in the source
// Disadvantages:
// * To many comment
//
// # Just show summary in comment
// Disadvantages:
// * Just will be some report that wouldn't be read
//

import (
	"github.com/egon12/ghr/coverage"
	"github.com/egon12/ghr/diff"
	"github.com/egon12/ghr/githubcommit"
	"github.com/egon12/ghr/path"
	"github.com/egon12/ghr/review"
)

// CoverageReviewer will report lines that not coverage in test
type CoverageReviewer interface {
	Read(coverProfilePath string) error
}

// NewCoverageReviewer create new CoverageReview
func NewCoverageReviewer(cs *githubcommit.Source, mc review.MultilineCommenter, rp review.Process) CoverageReviewer {
	return &coverageReviewer{cs, mc, rp, 0.7}
}

type coverageReviewer struct {
	commitSource       *githubcommit.Source
	multilineCommenter review.MultilineCommenter
	reviewProcess      ReviewProcess
	minimumForComment  float32
}

type coverProfile struct {
	file       string
	percentage float32
	ranges     []coverage.Range
}

func (c *coverageReviewer) Read(coverProfilePath string) error {
	commit := c.commitSource.GetCurrentCommit()
	if commit.Error() != nil {
		return commit.Error()
	}

	return c.Do(commit, coverProfilePath)
}

func (c *coverageReviewer) Do(com githubcommit.Commit, coverProfilePath string) error {
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

func (c *coverageReviewer) DoReview(com githubcommit.Commit, cv coverage.GoCoverageInGit, l diff.ListChanges) error {
	var filteredCoverProfile []coverProfile
	for _, f := range l.Files() {
		percent := cv.PercentageFile(f)
		if percent < 1 {
			ranges := cv.NotInCoverageLines(f)
			ranges = CombineRanges(ranges)
			filteredCoverProfile = append(filteredCoverProfile, coverProfile{
				file:       f,
				percentage: percent,
				ranges:     ranges,
			})
		}
	}

	return c.AddCoverageReview(com, filteredCoverProfile)
}

// AddCoverageReview will add review about coverage
func (c *coverageReviewer) AddCoverageReview(commit githubcommit.Commit, coverProfile []coverProfile) error {
	c.multilineCommenter.Start(commit)
	for _, cp := range coverProfile {
		err := c.AddSingleFileCoverageReview(cp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *coverageReviewer) AddSingleFileCoverageReview(cp coverProfile) error {
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
