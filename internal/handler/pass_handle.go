package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/senn404/bookmark-managent/internal/service"
)

type passwordHandler struct {
	svc service.Password
}

type Password interface {
	GenPass(c *gin.Context)
}

func NewPasswordHandler(svc service.Password) Password {
	return &passwordHandler{svc: svc}
}

func (h *passwordHandler) GenPass(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		c.String(http.StatusInternalServerError, "err")
		return
	}

	c.String(http.StatusOK, pass)
}
