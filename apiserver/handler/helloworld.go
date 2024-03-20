package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloWorld() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	}
}
