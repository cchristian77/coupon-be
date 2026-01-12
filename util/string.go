package util

import "strings"

func SanitizeString(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}
