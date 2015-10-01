package mailers

import (
	"taxiapp.com.br/app/models"
	"taxiapp.com.br/core"
)

func NewUserRegisteredEmail(user *models.User) *core.EmailTemplate {
	data := struct {
		Name string
	}{user.Name}

	return core.NewEmailTemplate(
		[]string{"layout/email.html", "user/registered.html"},
		user.Name,
		user.Email,
		"Bem vindo ao BrowTaxi",
		data,
	)
}

func NewResetPasswordEmail(user *models.User, token string) *core.EmailTemplate {
	data := struct {
		Name string
		Url  string
	}{user.Name, "http://taxiapp.com/auth/reset/" + token}

	return core.NewEmailTemplate(
		[]string{"layout/email.html", "user/reset_password.html"},
		user.Name,
		user.Email,
		"Para cadastrar uma nova senha siga as instruções",
		data,
	)
}

func NewChangedPasswordEmail(user *models.User) *core.EmailTemplate {
	data := struct {
		Name string
	}{user.Name}

	return core.NewEmailTemplate(
		[]string{"layout/email.html", "user/changed_password.html"},
		user.Name,
		user.Email,
		"Sua senha foi alterada",
		data,
	)
}
