package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExampleHandler struct {
	exampleService IExampleService
}

func NewExampleHandler(service IExampleService) ExampleHandler {
	return ExampleHandler{exampleService: service}
}

type IExampleService interface {
	Helloworld() string
}

func (h *ExampleHandler) HelloWorld() func(c *gin.Context) {
	return func(c *gin.Context) {
		s := h.exampleService.Helloworld()
		c.String(http.StatusOK, s)
	}
}
