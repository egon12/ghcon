package gotest

import (
	"bytes"
	"io"
	"io/ioutil"
)

func readStdOut(stdout io.Reader) (haveError bool, content string, err error) {
	// oc = outContent
	oc, err := ioutil.ReadAll(stdout)
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

		// what to do with this line, for now ignore it
		i = skipUntilNewline(oc, i)
	}

	if len(content) > 0 {
		haveError = true
		content = "Test Failed\n```\n" + content + "\n```\n"
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
