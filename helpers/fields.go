package helpers

import "regexp"

var countRegex = regexp.MustCompile(`(?i)COUNT\(.+\)`)

func IsCountOnly(fields []string) bool {
	if len(fields) != 1 {
		return false
	}
	return countRegex.MatchString(fields[0])
}
