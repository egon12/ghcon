package diff

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const maxIteration = 100000

// LineFinder interface for find linenumber in diff
type LineFinder interface {
	// Parse the diff content
	Parse([]byte) error

	// Find line number in diff content based on line number in the real file
	// isNew is to tell wether the line number is in original file or after
	// the file is edited
	Find(lineNumber int, isNew bool) (int, error)

	// StarHunk give the first line number that contain diff hunk (line start with
	// @@)
	StartHunk() (int, error)
}

type lineFinder struct {
	content []string
	hunks   []hunk
}

type hunk struct {
	lineNumber int
	oriStart   int
	oriLine    int
	newStart   int
	newLine    int
}

var hunkRegex = regexp.MustCompile(`^@@ -([0-9]+),([0-9]+) \+([0-9]+),([0-9]+) @@`)

// NewLineFinder Create new LineFinder
func NewLineFinder() LineFinder {
	return &lineFinder{}
}

func (d *lineFinder) Parse(b []byte) error {
	contents := string(b)
	d.fillContent(contents)
	d.fillHunks()
	if len(d.hunks) < 1 {
		return fmt.Errorf("Problably content is not a diff file (doesn't contain something like @@ -10,1 +10,1 @@)")
	}
	return nil
}

func (d *lineFinder) fillContent(contents string) {
	d.content = append([]string{""}, strings.Split(contents, "\n")...)
}

func (d *lineFinder) fillHunks() {
	for i, line := range d.content {
		if hunkRegex.MatchString(line) {
			d.parseHunk(i, line)
		}
	}
}

func (d *lineFinder) parseHunk(index int, line string) {
	submatch := hunkRegex.FindStringSubmatch(line)

	oriStart, _ := strconv.Atoi(submatch[1])
	oriLine, _ := strconv.Atoi(submatch[2])
	newStart, _ := strconv.Atoi(submatch[3])
	newLine, _ := strconv.Atoi(submatch[4])

	h := hunk{
		lineNumber: index,
		oriStart:   oriStart,
		oriLine:    oriLine,
		newStart:   newStart,
		newLine:    newLine,
	}
	d.hunks = append(d.hunks, h)
}

func (d *lineFinder) StartHunk() (int, error) {
	if len(d.hunks) < 1 {
		return 0, fmt.Errorf("This file doesn't have hunks or maybe forgot to Parse")
	}
	return d.hunks[0].lineNumber, nil
}

func (d *lineFinder) Find(n int, new bool) (int, error) {
	h, err := d.findHunk(n, new)
	if err != nil {
		return 0, err
	}

	lnFileCount := h.oriStart
	if new {
		lnFileCount = h.newStart
	}

	lnDiff := h.lineNumber + 1

	return d.countLineDiff(n, lnDiff, lnFileCount, new), nil
}

func (d *lineFinder) countLineDiff(lnInFile, lnDiff, lnFileCount int, new bool) int {
	var addRune byte = '-'
	var negRune byte = '+'
	if new {
		addRune = '+'
		negRune = '-'
	}

	for lnDiff < len(d.content) {
		firstChar := d.content[lnDiff][0]
		if lnInFile == lnFileCount && firstChar != negRune {
			break
		}

		switch firstChar {
		case ' ', addRune:
			lnFileCount++
		}
		lnDiff++

	}

	return lnDiff
}

func (d *lineFinder) findHunk(linenumber int, new bool) (hunk, error) {
	if new {
		return d.findHunkInNew(linenumber)
	}
	return d.findHunkInOri(linenumber)
}

func (d *lineFinder) findHunkInNew(linenumber int) (hunk, error) {
	for _, h := range d.hunks {
		if h.newStart <= linenumber && linenumber < h.newStart+h.newLine {
			return h, nil
		}
	}
	return hunk{}, fmt.Errorf("Cannot find line %d in new file from diff", linenumber)
}

func (d *lineFinder) findHunkInOri(linenumber int) (hunk, error) {
	for _, h := range d.hunks {
		if h.oriStart <= linenumber && linenumber < h.oriStart+h.oriLine {
			return h, nil
		}
	}
	return hunk{}, fmt.Errorf("Cannot find line %d in original file from diff", linenumber)
}
