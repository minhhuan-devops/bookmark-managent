package handler

import "github.com/gin-gonic/gin"

// genPassResponse represents the JSON response body for a successful
// password generation request.
type genPassResponse struct {
	Password string `json:"password"`
}

// errorResponse represents the standard JSON error response body
// returned when an API request fails.
type errorResponse struct {
	Error string `json:"error"`
}

// shortenURLResponse represents the JSON response body for a successful URL shortening request.
// Code is the generated short code that identifies the original URL.
// Message is a human-readable success message.
type shortenURLResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// responErr is a helper function that writes a JSON error response to the
// client with the given HTTP status code and error message.
func responErr(c *gin.Context, status int, msg string) {
	c.JSON(status, errorResponse{Error: msg})
}
