package validator

import (
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"
)

type (
	ValErrorSet map[string][]string

	Validator struct {
		Errors ValErrorSet
		Regex  *Regex
	}
)

var ValidatorMsg = newValMsg()

func NewValidator() Validator {
	return Validator{
		Errors: ValErrorSet(map[string][]string{}),
		Regex:  regexInstance(),
	}
}

func newValMsg() *valMsg {
	return &valMsg{
		RequiredErrMsg:   "required",
		MinLengthErrMsg:  "too short",
		MaxLengthErrMsg:  "too long",
		NotAllowedErrMsg: "not in allowed list",
		NotEmailErrMsg:   "not an email address",
		NoMatchErrMsg:    "confirmation does not match",
	}
}

type valMsg struct {
	RequiredErrMsg   string
	MinLengthErrMsg  string
	MaxLengthErrMsg  string
	NotAllowedErrMsg string
	NotEmailErrMsg   string
	NoMatchErrMsg    string
}

// ValidateRequired value.
func (v *Validator) ValidateRequired(val string, errMsg ...string) (ok bool) {
	val = strings.Trim(val, " ")
	return utf8.RuneCountInString(val) > 0
}

// ValidateMinLength value.
func (v *Validator) ValidateMinLength(val string, min int, errMsg ...string) (ok bool) {
	return utf8.RuneCountInString(val) >= min
}

// ValidateMaxLength value.
func (v *Validator) ValidateMaxLength(val string, max int) (ok bool) {
	return utf8.RuneCountInString(val) <= max
}

// ValidateEmail value.
func (v *Validator) ValidateEmail(val string) (ok bool) {
	return len(val) < 254 && v.Regex.EmailRegex.MatchString(val)
}

// ValidateConfirmation value.
func (v *Validator) ValidateConfirmation(val, confirmation string) (ok bool) {
	return val == confirmation
}

func (v *Validator) HasErrors() bool {
	return !v.IsValid()
}

func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}

func (es ValErrorSet) Add(field, msg string) {
	es[field] = append(es[field], msg)
}

func (es ValErrorSet) FieldErrors(field string) []string {
	return es[field]
}

func (es ValErrorSet) IsEmpty() bool {
	return len(es) < 1
}

type (
	Regex struct {
		MatchFirstCap *regexp.Regexp
		MatchAllCap   *regexp.Regexp
		EmailRegex    *regexp.Regexp
	}
)

var (
	lock = &sync.Mutex{}
	ri   *Regex
)

func regexInstance() *Regex {
	if ri == nil {
		lock.Lock()
		defer lock.Unlock()

		if ri == nil {
			ri = &Regex{
				MatchFirstCap: regexp.MustCompile("(.)([A-Z][a-z]+)"),
				MatchAllCap:   regexp.MustCompile("([a-z0-9])([A-Z])"),
				EmailRegex:    regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"),
			}

		}
	}

	return ri
}
