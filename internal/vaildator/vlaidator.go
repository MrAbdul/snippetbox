package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

// regex matching for email checks
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	FieldErrors map[string]string
}

// valid returns true if the fielderrors doesn't contain any entries
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// addFieldError adds an error message to the map as long as no entry already exists for the given key
func (v *Validator) AddFieldError(key, message string) {
	//we need to init the map if its not
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// checkfield adds an error message to the fielderrors map only if a validation check is not ok
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// our diffrent validators
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
