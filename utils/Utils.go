package utils

func ContainsOnly[T comparable](elems []T, elemToFind T) bool {
	for _, elem := range elems {
		if elem != elemToFind {
			return false
		}
	}
	return true
}

func GetAllIndicies[T comparable](elems []T, elemToFind T) []int {
	indicies := []int{}
	for idx, elem := range elems {
		if elem == elemToFind {
			indicies = append(indicies, idx)
		}
	}
	return indicies
}

func IsDigitsOnly(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
