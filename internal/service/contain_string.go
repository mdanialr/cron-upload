package service

// Contains tells whether b exist in a.
func Contains(a []string, b string) bool {
	for _, n := range a {
		if b == n {
			return true
		}
	}
	return false
}
