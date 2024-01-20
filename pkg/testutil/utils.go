package testutil

import (
	"reflect"
	"runtime"
	"strings"
)

// GetFunName returns the name of the function passed in.
func GetFunName(fn interface{}) string {
	if fn == nil {
		return ""
	}

	fullName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	parts := strings.Split(fullName, ".")
	return parts[len(parts)-1]
}
