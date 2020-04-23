package review

import (
	"fmt"

	"github.com/egon12/ghr/commit"
	"github.com/egon12/ghr/path"
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

func (p *processFacade) Comment(fileAndLine, comment string) error {
	sf := path.GetSourceFormatType(fileAndLine)
	switch sf {
	case path.FileAndLineNumber:
		filePath, line, err := path.ParseFileAndLine(fileAndLine)
		if err != nil {
			return fmt.Errorf("Cannot Parse File and Line: %v", err)
		}
		return p.reviewComment(filePath, line, comment)
	case path.FileAndRangeLine:
		filePath, from, to, err := path.ParseFileAndRangeLine(fileAndLine)
		if err != nil {
			return fmt.Errorf("Cannot Parse File and Range Line: %v", err)
		}
		return p.reviewMultilineComment(filePath, from, to, comment)
	}
	return fmt.Errorf("Cannot understand path and line: %s", fileAndLine)
}

func (p *processFacade) reviewComment(path string, lineNumber int, comment string) error {
	err := p.start()
	if err != nil {
		return err
	}

	return p.process.AddComment(path, lineNumber, comment)
}

func (p *processFacade) reviewMultilineComment(path string, fromLineNumber, toLineNumber int, comment string) error {
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
