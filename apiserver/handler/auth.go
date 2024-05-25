package handler

import (
	"net/http"
	"time"

	"protos/micro_auth"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"tools/errorx"
	"tools/oauth"
)

type AuthHandler struct {
	googleOauthConfig oauth2.Config
	grpcPools         GrpcConnectionPools
}

func NewAuthHandler(googleOauthConfig oauth2.Config, grpcPools GrpcConnectionPools) AuthHandler {
	return AuthHandler{googleOauthConfig: googleOauthConfig, grpcPools: grpcPools}
}

func (h *AuthHandler) GoogleOauthLoginHandler() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		url := h.googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
		ginCtx.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (h *AuthHandler) GoogleOauthCallbackHandler() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		if ginCtx.Query("code") == "" {
			ReturnClientBadRequestResponse(ginCtx, nil, "code is empty")
			return
		}

		token, err := h.googleOauthConfig.Exchange(
			ginCtx.Request.Context(),
			ginCtx.Query("code"),
		)
		if err != nil {
			ReturnClientBadRequestResponse(ginCtx, nil, "googleOauthConfig.Exchange err")
			return
		}
		if !token.Valid() {
			ReturnClientBadRequestResponse(ginCtx, nil, "token is invalid")
			return
		}

		userInfo, err := oauthx.GetGoogleUserInfo(token, 5*time.Second)
		if err != nil {
			ReturnServerInternalErrorResponse(ginCtx, nil, "oauthx.GetGoogleUserInfo err")
			return
		}
		conn, err := h.grpcPools.MicroAuthGrpcConnPool.GetConnFromPool()
		if errorx.IsErrorExist(err) {
			ReturnServerInternalErrorResponse(ginCtx, nil, "h.grpcPools.MicroAuthGrpcConnPool.GetConnFromPool err")
			return
		}
		if _, err = micro_auth.NewAuthServiceClient(conn).Login(
			ginCtx.Request.Context(),
			&micro_auth.LoginRequest{
				ID:            userInfo.ID,
				Email:         userInfo.Email,
				VerifiedEmail: userInfo.VerifiedEmail,
				Name:          userInfo.Name,
				GivenName:     userInfo.GivenName,
				FamilyName:    userInfo.FamilyName,
				Picture:       userInfo.Picture,
				Locale:        userInfo.Locale,
			},
		); errorx.IsErrorExist(err) {
			ReturnServerInternalErrorResponse(ginCtx, nil, "micro_auth.NewAuthServiceClient(conn).Login err")
			return
		}

		ReturnSuccessResponse(ginCtx, nil)
	}
}
