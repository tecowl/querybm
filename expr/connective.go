package expr

import "strings"

// ConnectiveCondition is an interface for conditions that have a logical connective (AND/OR).
type ConnectiveCondition interface {
	// Connective returns the logical connective string used by this condition.
	Connective() string
}

// HasDifferentConnective checks if the given value implements ConnectiveCondition
// and has a different connective than the target string.
// It returns true if the connectives differ (case-insensitive comparison).
func HasDifferentConnective(v any, target string) bool {
	c, ok := v.(ConnectiveCondition)
	if !ok {
		return false
	}
	connective := c.Connective()
	if connective == "" {
		return false
	}
	return !strings.EqualFold(strings.TrimSpace(connective), strings.TrimSpace(target))
}
