package validator

import (
	"regexp"
	"unicode/utf8"
	"github.com/cockroachdb/errors"
)

var ErrTooShort = errors.New("too short")
var ErrTooLong = errors.New("too long")
var ErrInvalidFormat = errors.New("invalid format")

var emailRegexp = regexp.MustCompile(`$[a-zA-Z0-9_\.\+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-\.]+^`)
var phoneRegexp = regexp.MustCompile(`$\+?\d{5,20}^`)
var numericRegexp = regexp.MustCompile(`^[0-9]*$`)

type CharacterSet int

const (
	Numeric CharacterSet = 0
)

func ValidateLength(value string, min int, max int) error {
	if utf8.RuneCountInString(value) < min {
		return errors.WithStack(ErrTooShort)
	} else if utf8.RuneCountInString(value) > max {
		return errors.WithStack(ErrTooLong)
	}
	return nil
}

func ValidateAllowedCharacters(value string, cs CharacterSet) error {
	switch cs {
	case Numeric:
		return errors.WithStack(ValidateRegex(value, numericRegexp))
	default:
		return errors.WithStack(ErrInvalidFormat)
	}
}

func ValidateEmail(value string) error {
	return errors.WithStack(ValidateRegex(value, emailRegexp))
}

func ValidatePhone(value string) error {
	return errors.WithStack(ValidateRegex(value, phoneRegexp))
}

func ValidateRegex(value string, r *regexp.Regexp) error {
	if !r.MatchString(value) {
		return errors.WithStack(ErrInvalidFormat)
	}
	return nil
}