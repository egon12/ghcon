package review

import (
	"fmt"

	"github.com/egon12/ghr/commit"
)

type processFacade struct {
	process            Process
	multilineCommenter MultilineCommenter
	source             CommitSource
}

func NewProcessFacade(p Process, m MultilineCommenter, c CommitSource) ProcessFacade {
	return &processFacade{
		process:            p,
		multilineCommenter: m,
		source:             c,
	}
}

func (p *processFacade) Comment(path string, lineNumber int, comment string) error {
	err := p.start()
	if err != nil {
		return err
	}

	return p.process.AddComment(path, lineNumber, comment)
}

func (p *processFacade) MultilineComment(path string, fromLineNumber, toLineNumber int, comment string) error {
	var (
		err    error
		commit commit.Commit
	)

	commit = p.source.GetCurrentCommit()
	if commit.Error() != nil {
		return err
	}

	if !commit.IsPR() {
		return fmt.Errorf("This current commit is not an PR, Cannot review Commit that not belong to any PR")
	}

	err = p.multilineCommenter.Start(commit)
	if err != nil {
		return err
	}

	return p.multilineCommenter.AddComment(path, fromLineNumber, toLineNumber, comment)
}

func (p *processFacade) Cancel() error {
	err := p.start()
	if err != nil {
		return err
	}

	return p.process.Cancel()
}

func (p *processFacade) Finish(lastComment string) error {
	err := p.start()
	if err != nil {
		return err
	}

	return p.process.Finish(lastComment)
}

func (p *processFacade) Approve(lastComment string) error {
	err := p.start()
	if err != nil {
		return err
	}

	return p.process.Approve(lastComment)
}

func (p *processFacade) RequestChanges(lastComment string) error {
	err := p.start()
	if err != nil {
		return err
	}

	return p.process.Finish(lastComment)
}

func (p *processFacade) start() error {
	var (
		err    error
		commit commit.Commit
	)

	commit = p.source.GetCurrentCommit()
	if commit.Error() != nil {
		return err
	}

	if !commit.IsPR() {
		return fmt.Errorf("This current commit is not an PR, Cannot review Commit that not belong to any PR")
	}

	err = p.process.Start(commit)
	if err != nil {
		return err
	}

	return nil
}
