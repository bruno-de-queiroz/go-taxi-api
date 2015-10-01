package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"github.com/creativelikeadog/go-taxi-api/app/forms"
	"github.com/creativelikeadog/go-taxi-api/app/models"
	"github.com/creativelikeadog/go-taxi-api/app/services"
	"github.com/creativelikeadog/go-taxi-api/core"
)

type DriverController struct {
	logger  *core.Logger
	service *services.DriverService
}

func (C *DriverController) Index(c *gin.Context) {

	page, err := strconv.Atoi(c.DefaultQuery("page", fmt.Sprintf("%d", DEFAULT_PAGE)))
	if err != nil {
		page = DEFAULT_PAGE
	}

	max, err := strconv.Atoi(c.DefaultQuery("max", fmt.Sprintf("%d", DEFAULT_PAGE_SIZE)))
	if err != nil {
		max = DEFAULT_PAGE_SIZE
	}

	offset := (page - 1) * max

	r, err := C.service.All(offset, max)
	if err != nil {
		C.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, r)
}

func (C *DriverController) Area(c *gin.Context) {

	var form forms.DriverInAreaForm

	err := c.Bind(&form)
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

	page := DEFAULT_PAGE
	max := NO_SIZE

	if form.Page != nil {
		page = *form.Page
	}

	if form.Max != nil {
		max = *form.Max
	}

	offset := (page - 1) * max

	if max == -1 {
		offset = (page - 1) * DEFAULT_PAGE_SIZE
	}

	r, err := C.service.InArea(form.GetArea(), offset, max)
	if err != nil {
		C.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, r)
}

func (C *DriverController) Create(c *gin.Context) {

	var form forms.DriverForm

	err := c.Bind(&form)
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

	if err := C.service.New(form); err != nil {
		if mgo.IsDup(err) {
			c.JSON(http.StatusConflict, []byte(nil))
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	c.JSON(http.StatusOK, []byte(nil))
}

func (C *DriverController) UpdateStatus(c *gin.Context) {

	var form forms.DriverStatusForm

	err := c.Bind(&form)
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

	driver := c.MustGet(DRIVER_ATTRIBUTE).(*models.Driver)

	if err := C.service.UpdateStatus(driver.ID, form); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, []byte(nil))

}

func (C *DriverController) Status(c *gin.Context) {
	m := c.MustGet(DRIVER_ATTRIBUTE).(*models.Driver)
	c.JSON(http.StatusOK, m)
}

func (C *DriverController) IsAreaParam(c *gin.Context) {
	//Bug do httprouter
	if c.Param(URI_PARAM) == "inArea" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func (C *DriverController) SetDriver(c *gin.Context) {
	//Bug do httprouter
	if c.Param(URI_PARAM) == "inArea" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	id := c.Param(URI_PARAM)

	if !bson.IsObjectIdHex(id) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	m, err := C.service.One(bson.ObjectIdHex(id))
	if err != nil {
		C.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Set(DRIVER_ATTRIBUTE, m)
	c.Next()
}

func NewDriverController(app *core.Application) *DriverController {
	return &DriverController{app.Logger, services.NewDriverService(app.Config.Database)}
}
