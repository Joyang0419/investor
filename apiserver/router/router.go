package router

import (
	"apiserver/handler"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	ExampleHandler handler.ExampleHandler
}

func NewHandlers(exampleHandler handler.ExampleHandler) Handlers {
	return Handlers{ExampleHandler: exampleHandler}
}

func NewGinRouter(handlers Handlers, middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.New()

	for idx := range middlewares {
		router.Use(middlewares[idx])
	}

	v1 := router.Group("/v1")
	{
		v1.GET("/helloworld", handlers.ExampleHandler.HelloWorld())
	}

	return router
}
