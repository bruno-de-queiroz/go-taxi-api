package forms
import (
	"github.com/creativelikeadog/taxiapp.com.br/app/exceptions"
)

type AuthForm struct {
	Email    *string `form:"email" json:"email"`
	Password *string `form:"password" json:"password"`
}

func (u *AuthForm) IsValid() (err error) {

	e := new(exceptions.ValidationException)

	if u.Email == nil || *u.Email == "" {
		e.Put("email", "is required")
	} else if !EMAIL_REGEX.MatchString(*u.Email) {
		e.Put("email", "is invalid")
	}

	if u.Password == nil || *u.Password == "" {
		e.Put("password", "is required")
	}

	if e.Size() != 0 {
		return e
	}

	return nil
}

type ResetForm struct {
	Email *string `form:"email" json:"email"`
}

func (u *ResetForm) IsValid() (err error) {

	e := new(exceptions.ValidationException)

	if u.Email == nil || *u.Email == "" {
		e.Put("email", "is required")
	} else if !EMAIL_REGEX.MatchString(*u.Email) {
		e.Put("email", "is invalid")
	}

	if e.Size() != 0 {
		return e
	}

	return nil
}

type PasswordForm struct {
	Password             *string `form:"password" json:"password"`
	PasswordConfirmation *string `form:"password_confirmation" json:"password_confirmation"`
}

func (u *PasswordForm) IsValid() (err error) {

	e := new(exceptions.ValidationException)

	if u.Password == nil || *u.Password == "" {
		e.Put("password", "is required")
	} else if len(*u.Password) < 8 {
		e.Put("password", "must have more than 8 caracters.")
	}

	if u.PasswordConfirmation == nil || *u.PasswordConfirmation == "" {
		e.Put("passwordConfirmation", "is required")
	}else if *u.Password != *u.PasswordConfirmation {
		e.Put("passwordConfirmation", "doesn't match.")
	}

	if e.Size() != 0 {
		return e
	}

	return nil
}