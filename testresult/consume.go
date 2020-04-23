package testresult

import (
	"bytes"
	"io/ioutil"
)

func Read(outPath, errPath string) (string, error) {
	haveError, content, err := ReadStdErr(errPath)
	if err != nil {
		return "", err
	}
	if haveError {
		return content, nil
	}

	haveError, content, err = ReadStdOut(outPath)
	if err != nil {
		return "", err
	}
	if haveError {
		return content, nil
	}

	return "", nil
}

func ReadStdErr(errPath string) (haveError bool, content string, err error) {
	errContent, err := ioutil.ReadFile(errPath)
	if err != nil {
		return
	}

	if len(errContent) > 0 {
		content = "Error when running test\n```\n" + string(errContent) + "\n```"
		haveError = true
	}

	return
}

func ReadStdOut(outPath string) (haveError bool, content string, err error) {
	// oc = outContent
	oc, err := ioutil.ReadFile(outPath)
	if err != nil {
		return
	}

	i := 0
	for i < len(oc)-1 {
		if bytes.Equal(oc[i:i+2], []byte("ok")) {
			// skip until "\n"
			i = skipUntilNewline(oc, i)
			continue
		}

		if bytes.Equal(oc[i:i+4], []byte("FAIL")) {
			startErr := i
			i = skipUntilNewline(oc, i)
			content += string(oc[startErr:i])
			continue
		}

		if bytes.Equal(oc[i:i+9], []byte("--- FAIL:")) {
			startErr := i
			endErr := i
			for endErr == startErr {
				i = skipUntilNewline(oc, i)
				if bytes.Equal(oc[i:i+4], []byte("FAIL")) {
					endErr = i
				}
			}
			content += string(oc[startErr:endErr])
			continue
		}

		panic(string(oc[i:len(oc)]))
	}

	if len(content) > 0 {
		haveError = true
		content = "Test Failed\n```\n" + content + "\n````"

	}

	return
}

func skipUntilNewline(oc []byte, i int) int {
	for i < len(oc) {
		if oc[i] == byte('\n') {
			return i + 1
		}
		i++
	}
	return i
}
