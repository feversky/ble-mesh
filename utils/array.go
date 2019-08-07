package utils

func Contains(arr []uint, e uint) bool {
	for _, a := range arr {
		if a == e {
			return true
		}
	}
	return false
}
