package graphql

import (
	"github.com/vektah/gqlparser/v2/gqlerror"

	"definition/response"
)

func gqlErr(err error, customCode response.TypeCustomCode, msg ...string) *gqlerror.Error {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}
	return &gqlerror.Error{
		Err:     err,
		Message: message,
		Extensions: map[string]interface{}{
			"code": customCode,
		},
	}
}

// Client Side

func ClientBadRequestErr(err error, msg ...string) *gqlerror.Error {
	return gqlErr(err, response.ClientBadRequest, msg...)
}

func ClientUnauthorizedErr(err error, msg ...string) *gqlerror.Error {
	return gqlErr(err, response.ClientUnauthorized, msg...)
}

// Server Side

func ServerInternalErr(err error, msg ...string) *gqlerror.Error {
	return gqlErr(err, response.ServerInternalError, msg...)
}
