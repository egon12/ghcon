package commit

type Commit interface {
	GetRepo() string
	GetHash() string
	IsPR() bool
	GetPRNumber() int
	GetPRID() string
	GetBaseRefName() string
	Error() error
}
