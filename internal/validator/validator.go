package validator

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

var ErrTooShort = errors.New("too short")
var ErrTooLong = errors.New("too long")


type CharacterSet int

const (
	Numeric CharacterSet = 0
	AlphaNumeric CharacterSet = 1
	Alphabetical CharacterSet = 2

)

func ValidateLength(value string, min int, max int) error {
	if utf8.RuneCountInString(value) < min {
		return ErrTooShort
	} else if utf8.RuneCountInString(value) > max {
		return ErrTooLong
	}
	return nil
}

func ValidateAllowedCharacters(value string, characterSet CharacterSet) error {
	switch characterSet {
	case Numeric:
	case AlphaNumeric:
	case Alphabetical:
	}
	return nil
}

func ValidateRegex(value string, regexp regexp.Regexp) error {
	return nil
}