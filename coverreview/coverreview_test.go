package coverreview

import (
	"testing"
)

func TestCoverageReviewer_Do(t *testing.T) {
	r := coverageReviewer{
		nil,
		&mockMultilineCommenter{},
		nil,
		0.7,
	}

	err := r.Do(nil, "")
	if err == nil {
		t.Error(err)
	}
	// TODO take out ListChanges Producer from this coverageReviewer
}

func TestCoverageReviewer_DoReview(t *testing.T) {
	r := coverageReviewer{
		nil,
		&mockMultilineCommenter{},
		nil,
		0.7,
	}

	err := r.DoReview(nil, &mockCoverage{}, &mockListChanges{})
	if err != nil {
		t.Error(err)
	}
}

// Ok, now how we will ad this one?
//go:generate mockery -name ReviewProcess -case snake -testonly -inpkg -keeptree
func TestCoverageReviewer_AddCoverageReview(t *testing.T) {
	/*
		rp := mocks.ReviewProcess{}
		rp.On("Start", "any")

		coverageReviewer{nil, nil, rp, 0.7}
	*/
}
