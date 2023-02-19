package internal

func CopySlice[T any](slice []T) []T {
	tmp := make([]T, len(slice))
	copy(tmp, slice)
	return tmp
}
