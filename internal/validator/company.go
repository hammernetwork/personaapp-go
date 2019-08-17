package validator

import (
	"github.com/cockroachdb/errors"
)

type Company struct {
	Name string
	Email string
	Phone string
	Password string
}

var ErrCompanyNameInvalid = errors.New("company_name is invalid")
var ErrCompanyEmailInvalid = errors.New("email is invalid")
var ErrCompanyPhoneInvalid = errors.New("phone is invalid")
var ErrCompanyPasswordInvalid = errors.New("password is invalid")

func ValidateCompany(c Company) error {
	if err := validateName(c.Name); err != nil {
		return errors.WithStack(ErrCompanyNameInvalid)
	}
	if err := validateEmail(c.Email); err != nil {
		return errors.WithStack(ErrCompanyEmailInvalid)
	}
	if err := validatePhone(c.Phone); err != nil {
		return errors.WithStack(ErrCompanyPhoneInvalid)
	}
	if err := validatePassword(c.Password); err != nil {
		return errors.WithStack(ErrCompanyPasswordInvalid)
	}
	return nil
}

func validateName(name string) error {
	if err := ValidateLength(name, 2,100); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func validateEmail(email string) error {
	if err := ValidateLength(email, 5,255); err != nil {
		return errors.WithStack(err)
	}
	if err := ValidateEmail(email); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func validatePhone(phone string) error {
	if err := validatePhone(phone); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func validatePassword(password string) error {
	if err := ValidateLength(password, 6, 30); err != nil {
		return errors.WithStack(err)
	}
	return nil
}