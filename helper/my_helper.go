package helper

func HasDuplicateString(arr []string) bool {
	counter := make(map[string]bool)
	for _, s := range arr {
		if counter[s] {
			return true
		}
		counter[s] = true
	}
	return false
}
