// Package handler provides HTTP request handlers for the Bookmark Management API.
package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/senn404/bookmark-managent/internal/service"
)

// ShortenURLHandler defines the interface for handling URL-shortening HTTP requests.
type ShortenURLHandler interface {
	// ShortenURL handles POST /shorten-url requests.
	// It accepts a URL and an expiration time, then returns a shortened code.
	ShortenURL(c *gin.Context)
	GetURL(c *gin.Context)
}

// shortenURLHandler is the concrete implementation of ShortenURLHandler.
// It delegates URL shortening logic to the ShortenURLService.
type shortenURLHandler struct {
	shortenURLService service.ShortenURLService
}

// NewShortenURLHandler creates a new ShortenURLHandler with the given ShortenURLService.
// It uses dependency injection to allow easy testing with mock services.
func NewShortenURLHandler(shortenURLService service.ShortenURLService) ShortenURLHandler {
	return &shortenURLHandler{
		shortenURLService: shortenURLService,
	}
}

// shortenURLRequest represents the JSON request body for the shorten URL endpoint.
// ExpTime is the desired expiration duration in minutes.
// URL is the original long URL to be shortened.
type shortenURLRequest struct {
	ExpTime int    `json:"exp_time" binding:"required,min=1"`
	URL     string `json:"url" binding:"required,url"`
}

// ShortenURL handles the POST /shorten-url endpoint.
// It parses and validates the request body, then delegates to the ShortenURLService.
// It returns a shortened URL code on success, or an error response on failure.
//
// @Summary Shorten URL
// @Tags URL Shortener
// @Accept json
// @Produce json
// @Param url body shortenURLRequest true "URL to shorten"
// @Success 200 {object} shortenURLRequest
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /links/shorten [post]
func (s *shortenURLHandler) ShortenURL(c *gin.Context) {
	input := &shortenURLRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		responErr(c, http.StatusBadRequest, "invaild input")
		return
	}

	key, err := s.shortenURLService.ShortenURL(c, input.URL, time.Duration(input.ExpTime)*time.Minute)
	if err != nil {
		log.Error().Err(err).Msg("Failed to shorten url")
		responErr(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, shortenURLResponse{
		Code:    key,
		Message: "Shorten URL generated successfully!",
	})
}

// GetURL handles GET /links/redirect/:code requests.
// It retrieves the original URL associated with the given code and redirects the client to it.
//
// @Summary Get Original URL
// @Tags URL Shortener
// @Produce json
// @Param code path string true "Shortened URL code"
// @Success 301 {string} string "Redirect to original URL"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /links/redirect/{code} [get]
func (s *shortenURLHandler) GetURL(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		responErr(c, http.StatusBadRequest, "invaild input")
		return
	}

	url, err := s.shortenURLService.GetURL(c, code)
	if err != nil {
		if errors.Is(err, service.ErrCodeNotExist) {
			responErr(c, http.StatusBadRequest, "code not exists")
			return
		}
		log.Error().Err(err).Msg("Failed to get url")
		responErr(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.Redirect(http.StatusMovedPermanently, url)
}
