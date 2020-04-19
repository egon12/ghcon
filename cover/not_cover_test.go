package cover

import "testing"

func TestGetNotCoverage(t *testing.T) {
	_, _ = GetNotCoverage("coverprofile.out")
}

func TestGetNotCoverage_CoverOut(t *testing.T) {
	cs, _ := GetNotCoverage("../cover.out")

	for _, nic := range cs {
		t.Logf("%v %v", nic.GetFile(), nic.GetRange())
	}
	//t.Error("what")
}
