package path

import (
	"fmt"
	"strconv"
	"strings"
)

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
