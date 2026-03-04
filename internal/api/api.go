package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/senn404/bookmark-managent/docs"
	"github.com/senn404/bookmark-managent/internal/handler"
	"github.com/senn404/bookmark-managent/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type api struct {
	app *gin.Engine
	cfg *Config
}

func (a *api) Start() error {
	return a.app.Run(fmt.Sprintf(":%s", a.cfg.AppPort))
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

func (a *api) registerEP() {
	passSvc := service.NewPassword()
	passHandler := handler.NewPasswordHandler(passSvc)
	a.app.GET("/gen-pass", passHandler.GenPass)
	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func New(cfg *Config) Engine {
	a := &api{
		app: gin.Default(),
		cfg: cfg,
	}
	a.registerEP()
	return a
}
