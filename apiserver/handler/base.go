package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"definition/response"
)

func resp(c *gin.Context, httpStatusCode int, customCode response.TypeCustomCode, data any, message ...string) {
	c.JSON(httpStatusCode,
		response.New(
			customCode,
			data,
			message...,
		),
	)
}

func SuccessResponse(
	c *gin.Context, data any, message ...string) {
	resp(c, http.StatusOK, response.Success, data, message...)
}

// ClientSide

func ClientBadRequestResponse(
	c *gin.Context, data any, message ...string) {
	resp(c, http.StatusBadRequest, response.ClientBadRequest, data, message...)
}

func ClientUnauthorizedResponse(
	c *gin.Context, data any, message ...string) {
	resp(c, http.StatusUnauthorized, response.ClientUnauthorized, data, message...)
}

// ServerSide

func ServerInternalErrorResponse(
	c *gin.Context, data any, message ...string) {
	resp(c, http.StatusInternalServerError, response.ServerInternalError, data, message...)
}
