package oauthx

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"tools/request"
)

/*
NewGoogleOauth creates a new oauth2.Config for Google OAuth2.
scopes: 訪問控制：授權服務器根據這些範圍來限制應用訪問用戶資料的範圍。只有用戶批准的範圍內的資料才可被應用訪問。
example:
https://www.googleapis.com/auth/userinfo.email：允許訪問用戶的郵件地址。
https://www.googleapis.com/auth/userinfo.profile：允許訪問用戶的基本個人資料信息。
https://www.googleapis.com/auth/drive：允許訪問和管理用戶的 Google 驅動器文件。
*/
func NewGoogleOauth(
	clientID string,
	clientSecret string,
	redirectURL string,
	scopes []GoogleOauthScope,
) oauth2.Config {
	return oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
	}
}

type GoogleOauthScope = string

const (
	ScopeForUserEmail   GoogleOauthScope = "https://www.googleapis.com/auth/userinfo.email"
	ScopeForUserProfile GoogleOauthScope = "https://www.googleapis.com/auth/userinfo.profile"
	ScopeForUserDrive   GoogleOauthScope = "https://www.googleapis.com/auth/drive"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`             // 用户在 Google 的唯一标识符
	Email         string `json:"email"`          // 用户的电子邮件地址
	VerifiedEmail bool   `json:"verified_email"` // 电子邮件是否经过验证
	Name          string `json:"name"`           // 用户的全名
	GivenName     string `json:"given_name"`     // 用户的名
	FamilyName    string `json:"family_name"`    // 用户的姓
	Picture       string `json:"picture"`        // 用户的头像图片 URL
	Locale        string `json:"locale"`         // 用户的区域设置
}

func GetGoogleUserInfo(token *oauth2.Token, timeout time.Duration) (userInfo *GoogleUserInfo, err error) {
	response, err := request.HttpRequest[GoogleUserInfo](
		"https://www.googleapis.com/oauth2/v2/userinfo",
		http.MethodGet,
		nil,
		timeout,
		map[string]string{
			"access_token": token.AccessToken,
		}, /* queryParams */
		nil, /* postBody */
	)
	if err != nil {
		return nil, fmt.Errorf("[GetGoogleUserInfo]request.HttpRequest error: %w", err)
	}

	return &response, nil
}
