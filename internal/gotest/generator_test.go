package gotest

import (
	"testing"

	"github.com/egon12/ghr/reviewing"
)

func TestGenerator_Generate(t *testing.T) {
	tg := NewGenerator()

	executor := NewMockExecutorWithFile("test-v0.1", "test_out", "test_err")

	stdout, stderr, _ := executor.Run(reviewing.PR{})

	review, err := tg.Generate(stdout, stderr, 0)
	if err != nil {
		t.Error(err)
	}

	if review == nil {
		t.Error("It should have a review")
	}

	if review.State != reviewing.Reject {
		t.Error("It should reject / request changes")
	}
}

func TestGenerate_Generate_TableTest(t *testing.T) {

	tests := []struct {
		name           string
		stdoutFilename string
		stderrFilename string
		messageOutLen  int
	}{
		{"Out some failed", "test_out", "test_err", 527},
		{"some stdout empty error", "test_out", "test_empty", 244},
		{"empty should not", "test_empty", "test_empty", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := NewGenerator()

			executor := NewMockExecutorWithFile("test-v0.1", tt.stdoutFilename, tt.stderrFilename)

			stdout, stderr, _ := executor.Run(reviewing.PR{})

			comments, _ := tg.Generate(stdout, stderr, 0)
			if tt.messageOutLen == 0 {
				if comments != nil {
					t.Errorf("Want 0 comments got %#v", comments)
				}
			} else if len(comments.Message) != tt.messageOutLen {
				t.Errorf("Diff len %d", len(comments.Message))
			}
		})
	}
}
