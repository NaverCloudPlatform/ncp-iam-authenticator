package utils

import (
	"strings"
)

func IsEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func ToPointer[T any](value T) *T {
	return &value
}
