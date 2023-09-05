package util

func AnyEmpty(strings ...string) bool {
	for _, s := range strings {
		if len(s) == 0 {
			return true
		}
	}
	return false
}
