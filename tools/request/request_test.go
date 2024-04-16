package request

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHttpRequest(t *testing.T) {
	// 模拟一个简单的 HTTP 服务器
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 根据请求路径返回不同的响应
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, `{"message":"success"}`)
		case "/forbidden":
			w.WriteHeader(http.StatusForbidden)
			_, _ = fmt.Fprintln(w, `{"message":"forbidden"}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	// 定义测试用例
	type response struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name                   string
		url                    string
		method                 string
		headers                map[string]string
		allowedHttpStatusCodes []int
		expectedResponse       response
		expectError            bool
	}{
		{
			name:                   "Success case",
			url:                    mockServer.URL + "/success",
			method:                 "GET",
			headers:                map[string]string{"Content-Type": "application/json"},
			allowedHttpStatusCodes: []int{http.StatusOK},
			expectedResponse: response{
				Message: "success",
			},
			expectError: false,
		},
		{
			name:                   "Error case",
			url:                    mockServer.URL + "/forbidden",
			method:                 "GET",
			headers:                map[string]string{"Content-Type": "application/json"},
			allowedHttpStatusCodes: []int{http.StatusOK},
			expectError:            true,
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp response
			resp, err := HttpRequest[response](tt.url, tt.method, tt.headers, 10*time.Second, tt.allowedHttpStatusCodes...)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, resp)
			}
		})
	}
}
