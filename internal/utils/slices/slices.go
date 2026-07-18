package slices

func Contains[T comparable](slice []T, a T) bool {
	for _, item := range slice {
		if item == a {
			return true
		}
	}
	return false
}
