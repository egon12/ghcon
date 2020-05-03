package coverreview

import (
	"reflect"
	"testing"

	"github.com/egon12/ghr/coverage"
)

// Ok let's pretend that this is some kind of interview code challange
// we are paint company that need to draw a line between two number.
// We got order by given two number. From and To. Unfortunately, sometimes
// the two number are overlap, and sometimes can be continued,
// but seperate into two range.
// Combine until we have minimum number of line that need to be draw
func TestCombineRanges(t *testing.T) {
	input := []coverage.Range{
		{1, 10},
		{11, 20},
		{3, 5},
	}

	want := []coverage.Range{{1, 20}}

	got := CombineRanges(input)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want\n%v\ngot %v", want, got)
	}
}

func TestCombineRanges_2(t *testing.T) {

	input := []coverage.Range{
		{34, 36},
		{41, 41},
		{44, 44},
		{36, 40},
		{41, 43},
		{47, 49},
		{56, 56},
		{59, 59},
		{49, 50},
		{50, 54},
		{56, 58},
		{62, 63},
		{68, 68},
		{63, 64},
		{64, 66},
		{100, 102},
		{108, 108},
		{111, 111},
		{102, 104},
		{104, 106},
		{108, 110},
		{114, 115},
		{121, 121},
		{115, 117},
		{117, 119},
	}

	want := []coverage.Range{
		{34, 44},
		{47, 54},
		{56, 59},
		{62, 66},
		{68, 68},
		{100, 106},
		{108, 111},
		{114, 119},
		{121, 121},
	}

	got := CombineRanges(input)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want\n%v\ngot %v", want, got)
	}

}

func TestNearStart(t *testing.T) {
	tests := []struct {
		name string
		r1   coverage.Range
		r2   coverage.Range
		want bool
	}{
		{
			name: "1,10 to 11,20",
			r1:   coverage.Range{From: 1, To: 10},
			r2:   coverage.Range{From: 11, To: 20},
			want: true,
		},
		{
			name: "1,9 to 11,20",
			r1:   coverage.Range{From: 1, To: 9},
			r2:   coverage.Range{From: 11, To: 20},
			want: false,
		},
		{
			name: "12,30 to 11,20",
			r1:   coverage.Range{From: 12, To: 30},
			r2:   coverage.Range{From: 11, To: 20},
			want: false,
		},
	}

	for _, tt := range tests {
		got := nearStart(tt.r1, tt.r2)
		if got != tt.want {
			t.Errorf("Want %v got %v", tt.want, got)
		}
	}
}

func TestInRange(t *testing.T) {
	r1 := coverage.Range{3, 5}
	r2 := coverage.Range{1, 20}
	want := true

	got := inRange(r1, r2)
	if inRange(r1, r2) != want {
		t.Errorf("Want %v got %v", want, got)
	}
}
