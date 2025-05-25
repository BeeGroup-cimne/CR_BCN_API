package utils

func IsValidEnum[T comparable](value T, validValues []T) bool {
	for _, v := range validValues {
		if v == value {
			return true
		}
	}
	return false
}
