package gotest

import (
	"bytes"
	"io"
	"os/exec"

	"github.com/egon12/ghr/reviewing"
)

// Executor runner for go test
type Executor struct{}

// ProvideTestExecutor for function to check
// if the Executor successfully implement reviewing.Executor
func ProvideTestExecutor() reviewing.Executor {
	return NewExecutor()
}

// NewTestExecutor Create new TestExecutor
func NewExecutor() *Executor {
	return &Executor{}
}

func (t *Executor) Name() string {
	return "test-v0.1"
}

func (t *Executor) Run(_ reviewing.PR, args ...string) (stdout io.Reader, stderr io.Reader, exitCode int) {

	args = append([]string{"test"}, args...)

	cmd := exec.Command("go", args...)
	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer

	err := cmd.Run()
	if err == nil {
		return stdoutBuffer, stderrBuffer, 0
	}

	exitError, ok := err.(*exec.ExitError)
	if ok {
		return stdoutBuffer, stderrBuffer, exitError.ExitCode()
	}

	return stdoutBuffer, stderrBuffer, 0
}
