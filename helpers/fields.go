// Package helpers provides utility functions for the querybm package.
package helpers

import "regexp"

var countRegex = regexp.MustCompile(`(?i)COUNT\(.+\)`)

// IsCountOnly checks if the fields contain only a single COUNT function.
// It returns true if there's exactly one field and it matches the COUNT(...) pattern.
func IsCountOnly(fields []string) bool {
	if len(fields) != 1 {
		return false
	}
	return countRegex.MatchString(fields[0])
}
