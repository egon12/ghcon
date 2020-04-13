package diff

import (
	"io/ioutil"
	"os"
	"testing"
)

var content []byte

var dlf LineFinder

func TestMain(m *testing.M) {
	f, _ := os.Open("thediff.diff")
	defer f.Close()

	content, _ = ioutil.ReadAll(f)
	dlf = &lineFinder{}
	_ = dlf.Parse(content)

	os.Exit(m.Run())
}

func TestLineFinder_StartHunk(t *testing.T) {
	got, err := dlf.StartHunk()
	if err != nil {
		t.Errorf("Unexpected Error %v", err)
	}

	if got != 5 {
		t.Errorf("Want 5 got %d", got)
	}
}

func TestLineFinder_ParseFailed(t *testing.T) {
	newLf := NewLineFinder()
	err := newLf.Parse([]byte("something for nothing"))
	if err == nil {
		t.Error("Expect Error")
	}
}

func TestLineFinder_Find(t *testing.T) {
	tests := []struct {
		name      string
		inputLine int
		inputNew  bool
		wantLine  int
		wantErr   bool
	}{
		{"Success in new firt line", 1, true, 7, false},
		{"Success in new second line", 2, true, 8, false},
		{"Success in Ori", 16, false, 16, false},
		{"Success in new 3 second hunk", 19, true, 22, false},
		{"Failed in new", 6, true, 0, true},
		{"Failed in ori", 5, false, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dlf.Find(tt.inputLine, tt.inputNew)
			if tt.wantErr == (err == nil) {
				t.Errorf("Unexpected Error %v", err)
			}

			if got != tt.wantLine {
				t.Errorf("Want %d got %d", tt.wantLine, got)
			}
		})
	}
}
