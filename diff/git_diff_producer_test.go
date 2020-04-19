package diff

import "testing"

func TestGitDiffProducer_Produce(t *testing.T) {
	gdp := &gitDiffProducer{}

	b, err := gdp.Produce(
		commit{
			hash:    "exp/second_try",
			baseref: "master",
		},
		"../ForDiff.md",
	)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)

	}

	if len(b) != 845 {
		t.Errorf("Maybe diff is not same as thediff.diff len %d", len(b))
	}
}

func TestGitDiffProducer_ListFile(t *testing.T) {
	gdp := &gitDiffProducer{}

	ls, err := gdp.ListFiles(
		commit{
			hash:    "exp/second_try",
			baseref: "master",
		},
	)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if ls[0] != "ForDiff.md" {
		t.Errorf("Want first item ForDiff.md got %v", ls)
	}

}

type commit struct {
	hash    string
	baseref string
}

func (c commit) GetHash() string        { return c.hash }
func (c commit) GetBaseRefName() string { return c.baseref }
