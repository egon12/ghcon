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

func TestLineFinder_Find_New(t *testing.T) {
	got, err := dlf.Find(19, true)
	if err != nil {
		t.Errorf("Unexpected Error %v", err)
	}

	if got != 22 {
		t.Errorf("Want 22 got %d", got)
	}
}

func TestLineFinder_Find_Ori(t *testing.T) {
	got, err := dlf.Find(16, false)
	if err != nil {
		t.Errorf("Unexpected Error %v", err)
	}

	if got != 16 {
		t.Errorf("Want 16 got %d", got)
	}
}

func TestLineFinder_Find_New_Failed(t *testing.T) {
	_, err := dlf.Find(6, true)
	if err == nil {
		t.Errorf("Expect Error")
	}
}

func TestLineFinder_Find_Ori_Failed(t *testing.T) {
	_, err := dlf.Find(5, false)
	if err == nil {
		t.Errorf("Unexpected Error %v", err)
	}
}

func TestLineFinder_ParseFailed(t *testing.T) {
	newLf := NewLineFinder()
	err := newLf.Parse([]byte("something for nothing"))
	if err == nil {
		t.Error("Expect Error")
	}
}

func TestLineFinder_Find_New_FirstLine(t *testing.T) {
	got, err := dlf.Find(1, true)
	if err != nil {
		t.Errorf("Unexpected Error %v", err)
	}

	if got != 7 {
		t.Errorf("Want 7 got %d", got)
	}
}
