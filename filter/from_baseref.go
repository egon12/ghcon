package filter

func FilterFiles(original, allowed []string) []string {
	var result []string
	for _, s := range original {
		if IsExistsIn(s, allowed) {
			result = append(result, s)
		}
	}
	return result
}

func IsExistsIn(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
