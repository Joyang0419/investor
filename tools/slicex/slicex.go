package slicex

import (
	"slices"
)

// IsEmpty 判斷切片為空
func IsEmpty[T any](s []T) bool {
	return len(s) == 0
}

// IsNotEmpty 判斷切片不為空
func IsNotEmpty[T any](s []T) bool {
	return len(s) != 0
}

// IsLengthFitExpected 判斷切片長度是否等於指定長度
func IsLengthFitExpected[T any](s []T, expected int) bool {
	return len(s) == expected
}

func GetLength[T any](s []T) int {
	return len(s)
}

func IsIdxInSlice[T any](s []T, idx int) bool {
	return idx >= 0 && idx < len(s)
}

func IsElementInSlice[S ~[]E, E comparable](s S, v E) bool {
	return slices.Contains(s, v)
}

func IsElementNotInSlice[S ~[]E, E comparable](s S, v E) bool {
	return !slices.Contains(s, v)
}

func RemoveDuplicateElement[T comparable](s []T) []T {
	elementMap := make(map[T]struct{})
	result := make([]T, 0, len(s))

	for _, elem := range s {
		if _, exists := elementMap[elem]; !exists {
			elementMap[elem] = struct{}{}
			result = append(result, elem)
		}
	}

	return result
}
