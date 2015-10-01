package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/http"
	"os"
	"path"
	"time"
)

const (
	CONFIG_PATH = "/config/"
	VIEW_PATH = "/views/"
)

type Application struct {
	*gin.Engine
	Config      *Config
	EmailSender *EmailSender
	Logger      *Logger
}

func NewApplication(p string, cfgs ...string) *Application {

	env := os.Getenv("APPLICATION_ENV")
	if env == "" {
		env = "development"
	}

	c, err := NewConfig(env, path.Join(p, CONFIG_PATH, "/environments/", env + ".yml"))
	if err != nil {
		panic(err)
	}

	for _, v := range cfgs {
		c.ExtendWithFile(path.Join(p, CONFIG_PATH, v))
	}

	log := NewLogger(c.Log)
	eh := NewEmailSender(c.Email, path.Join(p, "/app/views"))

	gin.SetMode(c.Mode)
	r := gin.New()

	r.Use(func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path

		ctx.Next()

		end := time.Now()
		latency := end.Sub(start)

		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()
		comment := ctx.Errors.ByType(gin.ErrorTypePrivate).String()

		log.Trace(statusCode,
			latency,
			clientIP,
			method,
			path,
			comment,
		)

	})

	r.Use(gin.Recovery())

	//TODO Verificar se os headers são suficientes
	r.Use(func(ctx *gin.Context) {

		ctx.Header("Cache-Control", "no-store")
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Max-Age", "86400")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, authorization, content-type")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length")
		ctx.Header("Access-Control-Allow-Methods", "POST, OPTIONS, DELETE, GET, PUT")
		ctx.Header("Connection", "keep-alive")
		ctx.Header("Date", time.Now().Format(time.RFC1123))
		ctx.Header("X-Request-Id", uuid.NewV4().String())

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		} else {
			ctx.Next()
		}
	})

	return &Application{r, c, eh, log}
}

func (a Application) Start() {
	a.Run(fmt.Sprintf("%s", a.Config.Host))
}
