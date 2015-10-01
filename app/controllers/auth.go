package controllers

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"github.com/creativelikeadog/taxiapp.com.br/app/forms"
	"github.com/creativelikeadog/taxiapp.com.br/app/services"
	"github.com/creativelikeadog/taxiapp.com.br/core"
)

type AuthController struct {
	logger  *core.Logger
	service *services.AuthService
}

func (C *AuthController) Authorize(c *gin.Context) {

	user, err := C.service.Authorize(c.Request)
	if err != nil {
		C.logger.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(CURRENT_USER_ATTRIBUTE, *user)
	c.Next()

}

func (C *AuthController) Reset(c *gin.Context) {

	var (
		form forms.ResetForm
		err error
	)

	err = c.Bind(&form)
	if err != nil {
		C.logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = form.IsValid()
	if err != nil {
		C.logger.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err = C.service.Reset(&form)
	if err != nil {
		C.logger.Error(err)
	}

	c.JSON(http.StatusOK, []byte(nil))
}

func (C *AuthController) NewPassword(c *gin.Context) {

	token := c.Param(TOKEN_PARAM)

	user, err := C.service.User(token)
	if err != nil {
		C.logger.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var (
		form forms.PasswordForm
	)

	err = c.Bind(&form)
	if err != nil {
		C.logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = form.IsValid()
	if err != nil {
		C.logger.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err = C.service.NewPassword(user, *form.Password)
	if err != nil {
		C.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, []byte(nil))
}

func (C *AuthController) Logout(c *gin.Context) {

	user, ok := c.Get(CURRENT_USER_ATTRIBUTE)

	if !ok {
		c.JSON(http.StatusOK, []byte(nil))
		return
	}

	u := user.(bson.ObjectId)

	err := C.service.Logout(u)
	if err != nil {
		C.logger.Error(err)
	}

	c.JSON(http.StatusOK, []byte(nil))
}

func (C *AuthController) Login(c *gin.Context) {

	var form forms.AuthForm

	err := c.Bind(&form)
	if err != nil {
		C.logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}

	err = form.IsValid()
	if err != nil {
		C.logger.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	token, err := C.service.Authenticate(&form)
	if err != nil {
		C.logger.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, struct {
		Token string `json:"token"`
	}{token})
}

func NewAuthController(app *core.Application) *AuthController {
	return &AuthController{app.Logger, services.NewAuthService(app)}
}
