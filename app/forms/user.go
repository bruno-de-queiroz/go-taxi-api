package forms

import (
	"taxiapp.com.br/app/exceptions"
)

type UserForm struct {
	Name                 *string `form:"name" json:"name"`
	Email                *string `form:"email" json:"email"`
	Password             *string `form:"password" json:"password"`
	PasswordConfirmation *string `form:"password_confirmation" json:"password_confirmation"`
}

func (u *UserForm) IsValid() (err error) {

	e := new(exceptions.ValidationException)

	if u.Email == nil || *u.Email == "" {
		e.Put("email", "is required")
	} else if !EMAIL_REGEX.MatchString(*u.Email) {
		e.Put("email", "is invalid.")
	}

	if u.Name == nil || *u.Name == "" {
		e.Put("name", "is required")
	}

	if u.Password == nil || *u.Password == "" {
		e.Put("password", "is required")
	} else if len(*u.Password) < 8 {
		e.Put("password", "must have more than 8 caracters.")
	}

	if u.PasswordConfirmation == nil || *u.PasswordConfirmation == "" {
		e.Put("passwordConfirmation", "is required")
	} else if *u.Password != *u.PasswordConfirmation {
		e.Put("passwords", "doesn't match.")
	}

	if e.Size() != 0 {
		return e
	}

	return nil
}
