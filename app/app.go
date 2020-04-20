package app

import (
	"github.com/egon12/ghr/commit"
	"github.com/egon12/ghr/review"
)

type App struct {
	ReviewProcess review.Process
	CommitSource  *commit.Source
}
