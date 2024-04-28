package api_response

import (
	"fmt"
)

type APIFormatResponse struct {
	Message    string `json:"message"`
	CustomCode int    `json:"customCode"`
	Data       any    `json:"data"`
}

func SetAPIFormatResponse(
	message string,
	customCode TypeCustomCode,
	data any,
) APIFormatResponse {
	return APIFormatResponse{
		Message:    fmt.Sprintf("%s: %s", getCustomCodeName(customCode), message),
		CustomCode: customCode,
		Data:       data,
	}
}

type TypeCustomCode = int

// 前三碼為 HTTP Status Code, 後四碼為自定義 Code
// 新增後，要記得去下面customCodeNames 新增名稱
const (
	SuccessQuery TypeCustomCode = 2000001

	ClientUnauthorized = 4010001

	ServerInternalError = 5000001
)

var customCodeNames = map[TypeCustomCode]string{
	SuccessQuery: "SuccessQuery",

	ClientUnauthorized: "ClientUnauthorized",

	ServerInternalError: "ServerInternalError",
}

func getCustomCodeName(code TypeCustomCode) string {
	customCodeName, keyExist := customCodeNames[code]
	if !keyExist {
		return "Unknown"
	}

	return customCodeName
}
