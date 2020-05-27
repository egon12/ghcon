package gotest

import (
	"io"
	"os"

	"github.com/egon12/ghr/reviewing"
)

type mockExecutor struct {
	name   string
	stdout io.Reader
	stderr io.Reader
}

func NewMockExecutorWithFile(name, outFile, errFile string) reviewing.Executor {
	stdout, err := os.Open(outFile)
	if err != nil {
		panic(err)
	}
	stderr, err := os.Open(errFile)
	if err != nil {
		panic(err)
	}
	return NewMockExecutor(name, stdout, stderr)
}

func NewMockExecutor(name string, stdout, stderr io.Reader) reviewing.Executor {
	return &mockExecutor{name, stdout, stderr}
}

func (m *mockExecutor) Name() string {
	return m.name
}

func (m *mockExecutor) Run(_ reviewing.PR, _ ...string) (stdout io.Reader, stderr io.Reader, exitCode int) {
	return m.stdout, m.stderr, 0
}
