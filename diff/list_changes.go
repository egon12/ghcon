package diff

func FromCommit(commit Commit) ListChanges {
	gd := GitDiffProducer{}
	files, err := gd.ListFiles(commit)
	if err != nil {
		panic(err)
	}
	return &listChanges{
		files:  files,
		commit: commit,
	}
}

type listChanges struct {
	files  []string
	commit Commit
}

func (l *listChanges) Files() []string {
	return l.files
}

func (l *listChanges) RangesInNew(filename string) []Range {
	gd := GitDiffProducer{}
	diff, err := gd.Produce(l.commit, filename)
	if err != nil {
		panic(err)
	}

	lf := NewLineFinder()
	lf.Parse(diff)
	oriLf := lf.(*lineFinder)

	var result []Range

	for _, h := range oriLf.hunks {
		result = append(result, Range{
			From: h.newStart,
			To:   h.newStart + h.newLine,
		})
	}

	return result
}

func (l *listChanges) RangesInOri(filename string) []Range {
	panic("not implemented") // TODO: Implement
}
