package lintreview

import (
	"github.com/egon12/ghr/commit"
	"github.com/egon12/ghr/diff"
	"github.com/egon12/ghr/filter"
	"github.com/egon12/ghr/lintresult"
	"github.com/egon12/ghr/path"
	"github.com/egon12/ghr/review"
)

type LintReviewer struct {
	reviewProcess review.Process
}

func NewLintReviewer(process review.Process) *LintReviewer {
	return &LintReviewer{process}
}

func (lr *LintReviewer) ReadAndReview(lintResultFile string, c commit.Commit) ([]lintresult.LintResult, error) {
	lintResult, err := lintresult.Read(lintResultFile)
	if err != nil {
		return nil, err
	}

	lc := diff.FromCommit(c)

	filesChanged := lc.Files()

	lr.reviewProcess.Start(c)
	for _, r := range lintResult {
		file, line, err := path.ParseFileAndLine(r.PathAndLine)
		if err != nil {
			return nil, err
		}

		if !filter.IsExistsIn(file, filesChanged) {
			continue
		}

		ranges := lc.RangesInNew(file)
		if !IsExistsWithinRanges(line, ranges) {
			continue
		}

		lr.reviewProcess.AddComment(file, line, r.Comment)
	}

	lr.reviewProcess.Finish("There are some automatic comment")

	return nil, nil
}

func IsExistsWithinRanges(line int, ranges []diff.Range) bool {
	for _, r := range ranges {
		if IsExistsWithinRange(line, r) {
			return true
		}
	}
	return false
}

func IsExistsWithinRange(line int, r diff.Range) bool {
	return r.From <= line && line <= r.To
}
