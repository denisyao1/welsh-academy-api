package util

func SliceHasNoDuplicate[T comparable](slice []T) bool {
	m := make(map[T]int)
	for _, v := range slice {
		_, exist := m[v]
		if exist {
			return false
		}
		m[v]++
	}
	return true
}

func Contains[T comparable](elm T, slice []T) bool {
	for _, v := range slice {
		if elm == v {
			return true
		}
	}
	return false
}
