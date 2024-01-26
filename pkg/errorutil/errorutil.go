package errorutil

import (
	"strings"
)

func ParseError(err error, msg string) bool {
	return strings.Contains(err.Error(), msg)
}
