package handler

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// Login by google oauth2
func Login(o oauth2.Config) gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func GoogleCallback(o oauth2.Config) gin.HandlerFunc {
	return func(c *gin.Context) {}
}
