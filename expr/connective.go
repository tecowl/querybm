package expr

import "strings"

type ConnectiveCondition interface {
	Connective() string
}

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
