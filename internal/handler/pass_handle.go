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

// @Summary Generate Password
// @Tags password
// @Produce json
// @Success 200 {object} genPassResponse
// @Failure 500 {object} errorResponse
// @Router /gen-pass [get]
func (h *passwordHandler) GenPass(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		responErr(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, genPassResponse{
		Password: pass,
	})
}
