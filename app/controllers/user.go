package controllers

import (
	"net/http"

	"github.com/creativelikeadog/go-taxi-api/app/forms"
	"github.com/creativelikeadog/go-taxi-api/app/services"
	"github.com/creativelikeadog/go-taxi-api/core"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	logger  *core.Logger
	service *services.UserService
}

func (C *UserController) Profile(c *gin.Context) {

	id := c.MustGet(CURRENT_USER_ATTRIBUTE).(bson.ObjectId)
	user, err := C.service.One(id)
	if err != nil {
		C.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)

}

func (C *UserController) Register(c *gin.Context) {

	var (
		form forms.UserForm
		err  error
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

	_, err = C.service.New(&form)
	if err != nil {
		C.logger.Error(err)
		if mgo.IsDup(err) {
			c.JSON(http.StatusConflict, []byte(nil))
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusCreated, nil)

}

func NewUserController(app *core.Application) *UserController {
	return &UserController{app.Logger, services.NewUserService(app)}
}
