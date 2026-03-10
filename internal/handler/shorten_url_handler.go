// Package handler provides HTTP request handlers for the Bookmark Management API.
package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/senn404/bookmark-managent/internal/service"
)

// ShortenURLHandler defines the interface for handling URL-shortening HTTP requests.
type ShortenURLHandler interface {
	// ShortenURL handles POST /shorten-url requests.
	// It accepts a URL and an expiration time, then returns a shortened code.
	ShortenURL(c *gin.Context)
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
	ExpTime int    `json:"exp_time"`
	URL     string `json:"url"`
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
// @Success 200 {object} shortenURLResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /shorten-url [post]
func (s *shortenURLHandler) ShortenURL(c *gin.Context) {
	input := &shortenURLRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		responErr(c, http.StatusBadRequest, "invaild input")
		return
	}

	key, err := s.shortenURLService.ShortenURL(c, input.URL, time.Duration(input.ExpTime)*time.Second)
	if err != nil {
		responErr(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, shortenURLResponse{
		Code:    key,
		Message: "Shorten URL generated successfully!",
	})
}
