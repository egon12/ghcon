package cover

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Range struct {
	From int
	To   int
}

type NotInCoverage interface {
	GetFile() string
	GetAbsPath() string
	GetRange() Range
}

type notInCoverage struct {
	Range
	file string
}

func (n *notInCoverage) GetFile() string {
	return n.file
}

func (n *notInCoverage) GetAbsPath() string {
	return os.Getenv("GOPATH") + "/src/" + n.file
}

func (n *notInCoverage) GetRange() Range {
	return n.Range
}

type CoverageStatus interface {
	GetFile() string
	GetRange() Range
	InCoverage() bool
}

func GetNotCoverage(filename string) ([]NotInCoverage, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	ns := strings.Split(string(b), "\n")

	var result []NotInCoverage

	for _, s := range ns {
		diff := strings.Split(s, ":")
		if len(diff) < 2 {
			continue
		}

		filename := diff[0]
		detail := diff[1]

		if filename == "mode" {
			continue
		}

		var lineStart int
		var columnStart int
		var lineEnd int
		var columnEnd int
		var unknown int
		var isCoverage int

		_, err := fmt.Sscanf(detail, "%d.%d,%d.%d %d %d", &lineStart, &columnStart, &lineEnd, &columnEnd, &unknown, &isCoverage)
		if err != nil {
			return nil, err
		}

		if isCoverage == 1 {
			continue
		}

		result = append(result, &notInCoverage{
			Range{lineStart, lineEnd},
			filename,
		})
	}
	return result, nil
}
