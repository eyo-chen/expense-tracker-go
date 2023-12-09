package validator

import "regexp"

var (
	EmailRX = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+\/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9-]+` + `(?:\.[a-zA-Z0-9-]+)*$`)
)

// Validator is a custom validator type which can hold a map of validation errors.
type Validator struct {
	Error map[string]string
}

// New creates a new Validator instance.
func New() *Validator {
	return &Validator{Error: make(map[string]string)}
}

// Valid returns true if the validator doesn't have any error, otherwise false.
func (v *Validator) Valid() bool {
	return len(v.Error) == 0
}

// Check adds an error message to the map of errors if the condition is false.
func (v *Validator) Check(valid bool, key, value string) {
	if valid {
		return
	}

	v.AddError(key, value)
}

// AddError adds an error message to the map of errors.
func (v *Validator) AddError(key, message string) {
	// Check if the key exists or not.
	if _, ok := v.Error[key]; ok {
		return
	}

	v.Error[key] = message
}

// Matches checks that a string value matches a specific regex pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
