package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"definition/response"
)

func returnResponse(c *gin.Context, httpStatusCode int, customCode response.TypeCustomCode, data any, message ...string) {
	c.JSON(httpStatusCode,
		response.New(
			customCode,
			data,
			message...,
		),
	)
}

func ReturnSuccessResponse(
	c *gin.Context, data any, message ...string) {
	returnResponse(c, http.StatusOK, response.Success, data, message...)
}

// ClientSide

func ClientBadRequestResponse(
	c *gin.Context, data any, message ...string) {
	returnResponse(c, http.StatusBadRequest, response.ClientBadRequest, data, message...)
}

func ClientUnauthorizedResponse(
	c *gin.Context, data any, message ...string) {
	returnResponse(c, http.StatusUnauthorized, response.ClientUnauthorized, data, message...)
}

// ServerSide

func ServerInternalErrorResponse(
	c *gin.Context, data any, message ...string) {
	returnResponse(c, http.StatusInternalServerError, response.ServerInternalError, data, message...)
}