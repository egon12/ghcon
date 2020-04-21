package path

import "testing"

func TestGitPath(t *testing.T) {
	gitPath, err := GitPath("./git_path.go")
	if err != nil {
		t.Errorf("Unexpected Error %v", err)
	}

	if gitPath != "path/git_path.go" {
		t.Errorf("Got %s", gitPath)
	}
}
