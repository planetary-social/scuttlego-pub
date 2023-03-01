package internal

func Pointer[T any](o T) *T {
	return &o
}
