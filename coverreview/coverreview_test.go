package coverreview

import (
	"testing"
)

func TestCoverageReviewer_Do(t *testing.T) {
	r := coverageReviewer{
		&mockMultilineCommenter{},
	}

	err := r.Do(nil, "")
	if err == nil {
		t.Error(err)
	}
	// TODO take out ListChanges Producer from this coverageReviewer
}

func TestCoverageReviewer_DoReview(t *testing.T) {
	r := coverageReviewer{
		&mockMultilineCommenter{},
	}

	err := r.DoReview(nil, &mockCoverage{}, &mockListChanges{})
	if err != nil {
		t.Error(err)
	}
}

func TestAddSingleCoverageReview(t *testing.T) {
}
