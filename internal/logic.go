package internal

import "github.com/egon12/ghr/reviewing"

type AutomaticReviewer struct {
	source   reviewing.GithubSource
	reviewer reviewing.Reviewer
}

func (a *AutomaticReviewer) Review(executor reviewing.Executor, generator reviewing.Generator, args ...string) error {
	pr, err := a.source.GetCurrentPR()
	if err != nil {
		return err
	}

	stdout, stderr, exitCode := executor.Run(pr, args...)

	comments, err := generator.Generate(stdout, stderr, exitCode)
	if err != nil {
		return err
	}

	// no comments and all is clear
	if len(comments) == 0 {
		return nil
	}

	err = a.reviewer.Start(pr)
	if err != nil {
		return err
	}

	for _, comment := range comments {
		if comment.CommentType == reviewing.CommentWithSource {
			err := a.reviewer.Comment(comment)
			if err != nil {
				return err
			}
		}

		if comment.CommentType == reviewing.CommentWithState {
			a.reviewer.Finish(comment)
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}
