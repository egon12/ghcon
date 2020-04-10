package diff

// GithubCommentPoisition is struct that can be used to find comment position
// when we review Github Pull Request
type GithubCommentPosition struct {
	dlf       LineFinder
	startHunk int
}

// CreateGithubCommentPostion it's input are the diffContent that can be get from
// git diff master...HEAD -- filename
// it will return error if the content are empty or not diff cotent
func NewGithubCommentPosition(diffContent []byte) (*GithubCommentPosition, error) {
	dlf := NewLineFinder()
	err := dlf.Parse(diffContent)
	if err != nil {
		return nil, err
	}

	startHunk, err := dlf.StartHunk()
	if err != nil {
		return nil, err
	}

	return &GithubCommentPosition{dlf, startHunk}, nil
}

// Find will get you comment position
func (g *GithubCommentPosition) Find(lineNumber int, isNew bool) (int, error) {
	ln, err := g.dlf.Find(lineNumber, isNew)
	if err != nil {
		return 0, err
	}

	return ln - g.startHunk, nil
}
