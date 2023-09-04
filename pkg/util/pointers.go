package util

func PointTo[T any](x T) *T {
	return &x
}
