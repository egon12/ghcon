package app

import (
	"github.com/egon12/ghr/coverreview"
	"github.com/egon12/ghr/review"
)

type App struct {
	ReviewProcess    review.ProcessFacade
	CoverageReviewer coverreview.CoverageReviewer
}
