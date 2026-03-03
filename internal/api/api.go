package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/senn404/bookmark-managent/internal/handler"
	"github.com/senn404/bookmark-managent/internal/service"
)

type api struct {
	app *gin.Engine
	cfg *Config
}

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func New(cfg *Config) Engine {
	a := &api{
		app: gin.New(),
		cfg: cfg,
	}
	a.registerEP()
	return a
}

func (a *api) Start() error {
	return a.app.Run(":" + a.cfg.AppPort)
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

func (a *api) registerEP() {
	passSvc := service.NewPassword()
	passHandler := handler.NewPasswordHandler(passSvc)
	a.app.GET("/gen-pass", passHandler.GenPass)
}
