package cover

import (
	"os"
	"testing"
)

func TestNotCoverExamples(t *testing.T) {
	NotCoverFuncExamples()
}

func TestFailed(t *testing.T) {
	failed := os.Getenv("FAILED")
	if failed == "Failed" {
		t.Error("Expected Failed Test")
	}
}
