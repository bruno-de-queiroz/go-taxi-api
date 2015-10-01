package services

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"github.com/creativelikeadog/taxiapp.com.br/app/exceptions"
	"github.com/creativelikeadog/taxiapp.com.br/app/forms"
	"github.com/creativelikeadog/taxiapp.com.br/app/mailers"
	"github.com/creativelikeadog/taxiapp.com.br/app/models"
	"github.com/creativelikeadog/taxiapp.com.br/core"
	"time"
)

type AuthService struct {
	logger *core.Logger
	user   *UserService
	config *core.TokenConfig
	sender *core.EmailSender
}

func (s *AuthService) hasExpired(t *jwt.Token) bool {
	timestamp := t.Claims["exp"]
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		return tm.Sub(time.Now()) <= 0
	}
	return true
}

func (s *AuthService) signFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(s.config.Secret), nil
}

func (s *AuthService) Authorize(r *http.Request) (id *bson.ObjectId, err error) {

	t, err := jwt.ParseFromRequest(r, s.signFunc)

	if err == nil && t.Valid {

		id := bson.ObjectIdHex(t.Claims["user"].(string))
		if s.hasExpired(t) {
			err = s.user.RemoveToken(id, ACCESS_TOKEN)
			if err != nil {
				return nil, err
			}

			return nil, &exceptions.TokenExpiredException{"Token is expired."}
		}

		u, err := s.user.One(id)
		if err != nil {
			return nil, err
		}

		if u.AccessToken == nil || *u.AccessToken != t.Raw {
			return nil, &exceptions.TokenNotFoundException{"Token not found"}
		}

		return &id, nil
	}

	return nil, err
}

func (s *AuthService) createToken(id bson.ObjectId) (token string, err error) {
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims["user"] = id.Hex()
	t.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(s.config.Expiration)).Unix()
	return t.SignedString([]byte(s.config.Secret))
}

func (s *AuthService) Reset(form *forms.ResetForm) (err error) {

	u, err := s.user.ByEmail(*form.Email)
	if err != nil {
		return err
	}

	if u.ResetToken != nil {
		t, err := jwt.Parse(*u.ResetToken, s.signFunc)
		if err == nil {
			if !s.hasExpired(t) {
				return nil
			}
		}
	}

	token, err := s.createToken(u.ID)
	if err != nil {
		return err
	}

	err = s.user.SaveToken(u.ID, RESET_TOKEN, token)
	if err != nil {
		return err
	}

	err = s.sender.Send(mailers.NewResetPasswordEmail(u, token))
	if err != nil {
		s.logger.Error(err)
	}

	return nil
}

func (s *AuthService) NewPassword(user bson.ObjectId, password string) (err error) {
	err = s.user.NewPassword(user, password)
	if err != nil {
		return err
	}

	u, err := s.user.One(user)
	if err != nil {
		return err
	}

	err = s.sender.Send(mailers.NewChangedPasswordEmail(u))
	if err != nil {
		s.logger.Error(err)
	}

	return nil
}

func (s *AuthService) User(token string) (u bson.ObjectId, err error) {
	t, err := jwt.Parse(token, s.signFunc)
	if err != nil {
		return bson.ObjectId(""), err
	}

	if s.hasExpired(t) {
		return bson.ObjectId(""), &exceptions.TokenExpiredException{"Token is expired."}
	}

	return bson.ObjectIdHex(t.Claims["user"].(string)), nil
}

func (s *AuthService) Logout(user bson.ObjectId) (err error) {
	return s.user.RemoveToken(user, ACCESS_TOKEN)
}

func (s *AuthService) Authenticate(form *forms.AuthForm) (token string, err error) {

	var u *models.User

	u, err = s.user.ByEmailAndPassword(*form.Email, *form.Password)
	if err != nil {
		return "", err
	}

	token, err = s.createToken(u.ID)
	if err != nil {
		return "", err
	}

	err = s.user.SaveToken(u.ID, ACCESS_TOKEN, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewAuthService(app *core.Application) *AuthService {
	return &AuthService{app.Logger, NewUserService(app), app.Config.Token, app.EmailSender}
}
