package request

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"tools/serialization"
)

type HttpRequestPostBody struct {
	ContentType string
	Body        any
}

func (body *HttpRequestPostBody) GetBody() (io.Reader, error) {
	switch body.ContentType {
	case "application/json":
		bodyBytes, err := serialization.JsonMarshal(body.Body)
		if err != nil {
			return nil, fmt.Errorf("[HttpRequestPostBody] serialization.JsonMarshal error: %w", err)
		}
		return bytes.NewReader(bodyBytes), nil
	case "application/x-www-form-urlencoded":
		data, ok := body.Body.(map[string]string)
		if !ok {
			return nil, fmt.Errorf("[HttpRequestPostBody] Error: Body type is not map[string]string for form encoding")
		}
		formData := url.Values{}
		for key, value := range data {
			formData.Set(key, value)
		}
		return strings.NewReader(formData.Encode()), nil
	default:
		return nil, fmt.Errorf("[HttpRequestPostBody]ContentType not supported: %v", body.ContentType)
	}
}

func HttpRequest[T any](
	urlStr,
	method string,
	headers map[string]string,
	timeout time.Duration,
	queryParams map[string]string,
	postBody *HttpRequestPostBody,
	allowedHttpStatusCodes ...int,
) (response T, err error) {
	// 创建一个 HTTP 请求
	if len(allowedHttpStatusCodes) == 0 {
		allowedHttpStatusCodes = []int{http.StatusOK}
	}
	// Parse URL and handle query parameters
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return response, fmt.Errorf("[HttpRequest]url.Parse error: %w", err)
	}

	// 處理queryParams
	query := parsedURL.Query()
	for key, value := range queryParams {
		query.Set(key, value)
	}
	parsedURL.RawQuery = query.Encode()

	// 處理postBody
	var bodyReader io.Reader
	if postBody != nil {
		bodyReader, err = postBody.GetBody()
		if err != nil {
			return response, fmt.Errorf("[HttpRequest]GetBody error: %w", err)
		}
	}

	req, err := http.NewRequest(method, parsedURL.String(), bodyReader)
	if err != nil {
		return response, fmt.Errorf("[HttpRequest]http.NewRequest err: %w", err)
	}

	// 添加请求头
	for key, value := range headers {
		req.Header.Set(key, value)
		// 防呆, 避免Headers 設定Content-Type，但PostBody 又設定了ContentType, 以PostBody 為主
		if postBody != nil {
			req.Header.Set("Content-Type", postBody.ContentType)
		}
	}

	// 发送请求
	client := &http.Client{Timeout: timeout}
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

	return serialization.JsonUnmarshal[T](bodyBytes)
}
