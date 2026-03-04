package handler

import "github.com/gin-gonic/gin"

type genPassResponse struct {
	Password string `json:"password"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func responErr(c *gin.Context, status int, msg string) {
	c.JSON(status, errorResponse{Error: msg})
}
