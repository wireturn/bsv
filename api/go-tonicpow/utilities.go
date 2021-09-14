package tonicpow

// isInList checks if string is known or not
func isInList(test string, list []string) bool {
	for _, a := range list {
		if test == a {
			return true
		}
	}
	return false
}
