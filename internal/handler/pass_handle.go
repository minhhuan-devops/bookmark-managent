package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/senn404/bookmark-managent/internal/service"
)

// passwordHandler is the concrete implementation of Password handler interface.
// It delegates password generation to the Password service.
type passwordHandler struct {
	svc service.Password
}

// Password defines the interface for handling password-related HTTP requests.
type Password interface {
	// GenPass handles GET /gen-pass requests and returns a randomly generated password.
	GenPass(c *gin.Context)
}

// NewPasswordHandler creates a new Password handler with the given Password
// service. It uses dependency injection to allow easy testing with mock services.
func NewPasswordHandler(svc service.Password) Password {
	return &passwordHandler{svc: svc}
}

// GenPass handles the GET /gen-pass endpoint.
// It generates a cryptographically secure random password and returns it
// as a JSON response. Returns HTTP 500 if password generation fails.
//
// @Summary Generate Password
// @Tags password
// @Produce json
// @Success 200 {object} genPassResponse
// @Failure 500 {object} errorResponse
// @Router /gen-pass [get]
func (h *passwordHandler) GenPass(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate password")
		responErr(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, genPassResponse{
		Password: pass,
	})
}
