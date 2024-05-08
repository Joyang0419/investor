package slicex

// IsEmpty 判斷切片為空
func IsEmpty[T any](s []T) bool {
	return len(s) == 0
}

// IsNotEmpty 判斷切片不為空
func IsNotEmpty[T any](s []T) bool {
	return len(s) != 0
}

// CheckLengthFitExpected 判斷切片長度是否等於指定長度
func CheckLengthFitExpected[T any](s []T, expected int) bool {
	return len(s) == expected
}

func GetLength[T any](s []T) int {
	return len(s)
}

func CheckIdxInSlice[T any](s []T, idx int) bool {
	return idx >= 0 && idx < len(s)
}
