package path

import (
	"fmt"
	"strconv"
	"strings"
)

type SourceFormat int

const (
	FilePathOnly SourceFormat = iota + 1
	FileAndLineNumber
	FileAndLineWithColumn
	FileAndRangeLine
	FileAndRangeLineWithColumn

	Unknwon SourceFormat = 99
)

func GetSourceFormatType(input string) SourceFormat {
	pnl := strings.Split(input, ":")
	if len(pnl) == 1 {
		return FilePathOnly
	}

	if len(pnl) == 2 {
		if !strings.Contains(pnl[1], ".") {
			return FileAndLineNumber
		}
		return FileAndLineWithColumn
	}

	if len(pnl) == 3 {
		if !strings.Contains(pnl[1], ".") {
			return FileAndRangeLine
		}
		return FileAndRangeLineWithColumn
	}

	return Unknwon
}

func ParseFileAndRangeLine(input string) (path string, from, to int, err error) {
	pnl := strings.Split(input, ":")
	if len(pnl) != 3 {
		err = fmt.Errorf("Format of the path and line number should be path/file.go:10:13")
		return
	}

	path = pnl[0]

	from, err = strconv.Atoi(pnl[1])
	if err != nil {
		return
	}

	to, err = strconv.Atoi(pnl[2])
	return
}

func ParseFileAndLine(pathAndLine string) (path string, line int, err error) {
	pnl := strings.Split(pathAndLine, ":")
	if len(pnl) != 2 {
		err = fmt.Errorf("Format of the path and line number should be path/file.go:10")
		return
	}

	path = pnl[0]
	lineStr := pnl[1]

	line, err = strconv.Atoi(lineStr)

	return path, line, err
}
