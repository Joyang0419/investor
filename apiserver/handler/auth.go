package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"tools/oauthx"
)

func GoogleOauthLoginHandler(googleOauthConfig oauth2.Config) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
		ginCtx.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func GoogleOauthCallbackHandler(googleOauthConfig oauth2.Config) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		if ginCtx.Query("code") == "" {
			ginCtx.JSON(http.StatusBadRequest, gin.H{
				"error": "code is required",
			})
			return
		}

		token, err := googleOauthConfig.Exchange(
			ginCtx.Request.Context(),
			ginCtx.Query("code"))
		if err != nil {
			_ = ginCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to exchange token: %v", err))
		}
		if !token.Valid() {
			_ = ginCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("invalid token: %v", token))
		}

		userInfo, err := oauthx.GetGoogleUserInfo(token, 5*time.Second)
		if err != nil {
			_ = ginCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user info: %v", err))
		}
		_ = userInfo

		ginCtx.JSON(http.StatusOK, gin.H{
			"status": "You are logged in",
		})
	}
}
