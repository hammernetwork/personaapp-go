package models

import "persona/forms"

type User struct {
}

func (u User) Signup(userPayload forms.UserSignup) error {
	return nil
}
