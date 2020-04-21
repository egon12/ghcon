package coverreview

import (
	"github.com/egon12/ghr/commit"
	"github.com/egon12/ghr/cover"
	"github.com/egon12/ghr/diff"
	"github.com/egon12/ghr/path"
	"github.com/egon12/ghr/review"
)

type coverageReviewer struct {
	gitDiffProducer    diff.GitDiffProducer
	multilineCommenter review.MultilineCommenter
}

func (c *coverageReviewer) Do(commit commit.Commit, coverProfilePath string) error {

	listNotInCover, err := cover.GetNotCoverage(coverProfilePath)
	if err != nil {
		return err
	}

	listFileChanges, err := c.gitDiffProducer.ListFiles(commit)
	if err != nil {
		return err
	}

	var filteredNotInCover []cover.NotInCoverage
	for _, nic := range listNotInCover {
		for _, filePath := range listFileChanges {
			gitFile, _ := path.GitPath(nic.GetAbsPath())
			if gitFile == filePath {
				filteredNotInCover = append(filteredNotInCover, nic)
			}
		}
	}

	return c.AddCoverageReview(commit, filteredNotInCover)

}

func (c *coverageReviewer) filterNotInCoverage(input []cover.NotInCoverage, listCommitFiles []string) []cover.NotInCoverage {
	var filteredNotInCover []cover.NotInCoverage
	for _, nic := range input {
		for _, filePath := range listCommitFiles {
			gitFile, _ := path.GitPath(nic.GetAbsPath())
			if gitFile == filePath {
				filteredNotInCover = append(filteredNotInCover, nic)
			}
		}
	}
	return filteredNotInCover
}

func (c *coverageReviewer) AddCoverageReview(commit commit.Commit, nics []cover.NotInCoverage) error {
	c.multilineCommenter.Start(commit)
	for _, nic := range nics {
		err := c.AddSingleCoverageReview(nic)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *coverageReviewer) AddSingleCoverageReview(nic cover.NotInCoverage) error {
	lineRange := nic.GetRange()
	gitPath, _ := path.GitPath(nic.GetAbsPath())
	return c.multilineCommenter.AddComment(gitPath, lineRange.From, lineRange.To, "Not In Coverage")
}
