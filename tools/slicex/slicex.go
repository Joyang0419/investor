package slicex

// IsEmpty 判斷切片為空
func IsEmpty[T any](s []T) bool {
	return len(s) == 0
}

// IsNotEmpty 判斷切片不為空
func IsNotEmpty[T any](s []T) bool {
	return len(s) != 0
}
