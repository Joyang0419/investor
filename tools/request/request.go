package request

import (
	"fmt"
	"io"
	"net/http"
	"slices"

	"tools/serialization"
)

var client = &http.Client{}

func HttpRequest[T any](url, method string, headers map[string]string, allowedHttpStatusCodes ...int) (response T, err error) {
	// 创建一个 HTTP 请求
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return response, fmt.Errorf("[HttpRequest]http.NewRequest err: %w", err)
	}

	// 添加请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("[HttpRequest]client.Do err: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 检查响应状态码
	if !slices.Contains(allowedHttpStatusCodes, resp.StatusCode) {
		return response, fmt.Errorf("[HttpRequest] resp.StatusCode not in allowedHttpStatusCodes, resp.StatusCode:%v, allowedHttpStatusCodes: %v", resp.StatusCode, allowedHttpStatusCodes)
	}

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("[HttpRequest]io.ReadAll err: %w", err)
	}

	return serialization.JsonUnmarshal[T](string(bodyBytes))
}
