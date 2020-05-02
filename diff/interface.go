package diff

type Range struct {
	From int
	To   int
}

type ListChanges interface {
	Files() []string
	RangesInNew(filename string) []Range
	RangesInOri(filename string) []Range
}
