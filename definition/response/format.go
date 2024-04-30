package response

import (
	"fmt"
)

type Response struct {
	Message    string `json:"message"`
	CustomCode int    `json:"customCode"`
	Data       any    `json:"data"`
}

func New(
	message string,
	customCode TypeCustomCode,
	data any,
) Response {
	return Response{
		Message:    fmt.Sprintf("%s: %s", getCustomCodeName(customCode), message),
		CustomCode: customCode,
		Data:       data,
	}
}

type TypeCustomCode = int

// 前三碼為 HTTP Status Code, 後四碼為自定義 Code
// 新增後，要記得去下面customCodeNames 新增名稱
const (
	Success TypeCustomCode = 2000000

	ClientUnauthorized = 4010000
	// ClientConflict 資源已存在
	ClientConflict = 4090000

	ServerInternalError = 5000000
)

var customCodeNames = map[TypeCustomCode]string{
	Success: "Success",

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
