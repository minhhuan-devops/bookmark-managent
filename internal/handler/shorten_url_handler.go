package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/senn404/bookmark-managent/internal/service"
)

type ShortenURLHandler interface {
	ShortenURL(c *gin.Context)
}

type shortenURLHandler struct {
	shortenURLService service.ShortenURLService
}

func NewShortenURLHandler(shortenURLService service.ShortenURLService) ShortenURLHandler {
	return &shortenURLHandler{
		shortenURLService: shortenURLService,
	}
}

type shortenURLRequest struct {
	ExpTime int    `json:"exp_time"`
	URL     string `json:"url"`
}

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

	key, err := s.shortenURLService.ShortenURL(c, input.URL, time.Duration(input.ExpTime)*time.Minute)
	if err != nil {
		responErr(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, shortenURLResponse{
		Code: key,
		Message: "Shorten URL generated successfully!",
	})
}
