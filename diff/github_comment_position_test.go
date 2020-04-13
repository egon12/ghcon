package diff

import "testing"

func TestGithubCommentPosition(t *testing.T) {

	tests := []struct {
		name      string
		inputLine int
		inputNew  bool
		wantLine  int
		wantErr   bool
	}{
		{"first line", 1, true, 2, false},
		{"second line", 2, true, 3, false},
		{"second hunk ori", 16, false, 11, false},
		{"second hunk in new", 19, true, 17, false},
	}

	c := commit{
		hash:    "exp/third_try",
		baseref: "master",
	}

	g, _ := NewGithubCommentPosition(c, "../ForDiff.md")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := g.Find(tt.inputLine, tt.inputNew)

			if got != tt.wantLine {
				t.Errorf("Want %d got %d", tt.wantLine, got)
			}
		})
	}
}
