package utils

func Contains(value []string, v string) bool {
	for _, val := range value {
		if val == v {
			return true
		}
	}
	return false
}
