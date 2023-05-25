package kc

type GenericType interface {
	int64 | float64 | string
}

func contains[T GenericType](array []T, value T) bool {
	for _, d := range array {
		if d == value {
			return true
		}
	}
	return false
}

func findIndex[T GenericType](array []T, value T) int {
	for i, d := range array {
		if d == value {
			return i
		}
	}
	return -1
}

func remove[T GenericType](array []T, index int) []T {
	return append(array[:index], array[index+1:]...)
}
