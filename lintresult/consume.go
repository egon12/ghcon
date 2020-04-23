package lintresult

import (
	"bytes"
	"io/ioutil"
)

type LintResult struct {
	PathAndLine string
	Comment     string
}

func Read(filename string) ([]LintResult, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	content := bytes.Split(b, []byte("\n"))

	content = filterDetail(content)

	return mapToLintResult(content), nil
}

func filterDetail(input [][]byte) [][]byte {
	var result [][]byte
	for _, b := range input {
		if len(b) < 1 {
			continue
		}

		if b[0] == byte('\t') {
			continue
		}

		result = append(result, b)
	}

	return result
}

func mapToLintResult(input [][]byte) []LintResult {
	var result []LintResult
	for _, b := range input {
		bs := bytes.SplitN(b, []byte(" "), 2)
		pathLineColumn := bs[0]
		pathLine := pathLineColumn[:len(pathLineColumn)-3]

		result = append(result, LintResult{
			PathAndLine: string(pathLine),
			Comment:     string(bs[1]),
		})
	}
	return result
}
