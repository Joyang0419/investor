package router

import (
	"github.com/gin-gonic/gin"

	"apiserver/handler"
)

type Handler struct {
	AuthHandler    handler.AuthHandler
	GraphqlHandler handler.GraphqlHandler
}

func NewGinRouter(
	middlewares []gin.HandlerFunc,
	handler Handler,
) *gin.Engine {
	router := gin.New()

	for idx := range middlewares {
		router.Use(middlewares[idx])
	}

	router.POST("/query",
		handler.AuthHandler.ValidateToken(),
		handler.GraphqlHandler.HandleGraphql(),
	)

	// 尚未進到系統前，都使用Restful api(example: login, callback, ...)
	router.GET("/auth/google/login", handler.AuthHandler.GoogleOauthLoginHandler())
	router.GET("/auth/google/callback", handler.AuthHandler.GoogleOauthCallbackHandler())

	return router
}
