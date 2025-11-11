package utils

import "strings"

func StandardError(err error) string {
	errorMessages := strings.Split(err.Error(), "]: ")
	return errorMessages[len(errorMessages)-1]
}
