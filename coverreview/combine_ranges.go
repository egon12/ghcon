package coverreview

import "github.com/egon12/ghr/coverage"

// CombineRanges will combine ranges that near and overlap
func CombineRanges(input []coverage.Range) []coverage.Range {
	result := combineRanges(input)
	return combineRanges(result)
}

func combineRanges(input []coverage.Range) []coverage.Range {
	var result []coverage.Range
	for _, r1 := range input {
		var added = false
		for i, r2 := range result {
			if inRange(r1, r2) {
				added = true
				break
			}

			if nearStart(r1, r2) {
				added = true
				result[i].From = r1.From
			}

			if nearEnd(r1, r2) {
				added = true
				result[i].To = r1.To
			}

			if added {
				break
			}

			if startOverlap(r1, r2) {
				added = true
				result[i].To = r1.To
			}

			if endOverlap(r1, r2) {
				added = true
				result[i].From = r1.From
			}

			if added {
				break
			}
		}
		if !added {
			result = append(result, r1)
		}
	}
	return result
}

func inRange(r1 coverage.Range, r2 coverage.Range) bool {
	return r1.From >= r2.From && r1.To <= r2.To
}

func nearStart(r1 coverage.Range, r2 coverage.Range) bool {
	return r1.To == r2.From-1 || r1.To == r2.From
}

func nearEnd(r1 coverage.Range, r2 coverage.Range) bool {
	return r1.From == r2.To+1 || r1.From == r2.To
}

func startOverlap(r1 coverage.Range, r2 coverage.Range) bool {
	return r2.From <= r1.From && r1.From <= r2.To && r1.To > r2.To
}

func endOverlap(r1 coverage.Range, r2 coverage.Range) bool {
	return r2.From <= r1.To && r1.To <= r2.To && r1.From < r2.From
}
