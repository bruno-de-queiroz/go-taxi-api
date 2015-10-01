package main

import (
	"path/filepath"
	"github.com/creativelikeadog/taxiapp.com.br/app/controllers"
	"github.com/creativelikeadog/taxiapp.com.br/core"
)

func main() {

	dir, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}

	a := core.NewApplication(dir)

	auth := controllers.NewAuthController(a)

	a.POST("/login", auth.Login)
	a.DELETE("/logout", auth.Authorize, auth.Logout)

	a.POST("/reset", auth.Reset)
	a.POST("/reset/:" + controllers.TOKEN_PARAM, auth.NewPassword)

	uc := controllers.NewUserController(a)
	a.POST("/register", uc.Register)
	a.GET("/profile", auth.Authorize, uc.Profile)

	dc := controllers.NewDriverController(a)
	drivers := a.Group("/drivers", auth.Authorize)
	{
		drivers.GET("/", dc.Index)
		drivers.POST("/", dc.Create)
		param := drivers.Group("/:" + controllers.URI_PARAM)
		{
			//Bug do httprouter
			param.GET("", dc.IsAreaParam, dc.Area)
			param.GET("/status", dc.SetDriver, dc.Status)
			param.PUT("/status", dc.SetDriver, dc.UpdateStatus)
		}
	}

	a.Start()
}
