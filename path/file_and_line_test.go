package path

import "testing"

func TestParseFileAndLine(t *testing.T) {
	fileAndLine := "path/file_and_line.go:10"

	path, line, err := ParseFileAndLine(fileAndLine)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if path != "path/file_and_line.go" {
		t.Errorf("Wrong path got %s", path)
	}

	if line != 10 {
		t.Errorf("Wrong line number got %d", line)
	}
}
