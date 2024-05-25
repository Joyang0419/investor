package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"protos/micro_auth"
	"tools/stringx"
	"tools/timex"

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
			// 用這個舉例好了, 你會知道ReturnClientBadRequestResponse, 走完左邊那個方法，就會return掉了嗎?
			ClientBadRequestResponse(ginCtx, nil, "code is empty")
			return
		}

		token, err := h.googleOauthConfig.Exchange(
			ginCtx.Request.Context(),
			ginCtx.Query("code"),
		)
		if err != nil {
			ClientBadRequestResponse(ginCtx, nil, "googleOauthConfig.Exchange err")
			return
		}
		if !token.Valid() {
			ClientBadRequestResponse(ginCtx, nil, "token is invalid")
			return
		}

		userInfo, err := oauthx.GetGoogleUserInfo(token, 5*time.Second)
		if err != nil {
			ServerInternalErrorResponse(ginCtx, nil, "oauthx.GetGoogleUserInfo err")
			return
		}
		conn, err := h.grpcPools.MicroAuthGrpcConnPool.GetConnFromPool()
		if errorx.IsErrorExist(err) {
			ServerInternalErrorResponse(ginCtx, nil, "h.grpcPools.MicroAuthGrpcConnPool.GetConnFromPool err")
			return
		}
		grpcResp, err := micro_auth.NewAuthServiceClient(conn).Login(
			ginCtx.Request.Context(),
			&micro_auth.LoginRequest{
				ID:                 userInfo.ID,
				Email:              userInfo.Email,
				VerifiedEmail:      userInfo.VerifiedEmail,
				Name:               userInfo.Name,
				GivenName:          userInfo.GivenName,
				FamilyName:         userInfo.FamilyName,
				Picture:            userInfo.Picture,
				Locale:             userInfo.Locale,
				LastLoginTimestamp: timex.GetCurrentTimestampSeconds(),
			},
		)
		if errorx.IsErrorExist(err) {
			ServerInternalErrorResponse(ginCtx, nil, "micro_auth.NewAuthServiceClient(conn).Login err")
			return
		}

		ReturnSuccessResponse(ginCtx, gin.H{"token": grpcResp.Token})
	}
}

const TokenInfoCtxKey = "TokenInfo"

type TokenInfo struct {
	ID    string
	Email string
}

func (h *AuthHandler) ValidateToken() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		authorization := ginCtx.GetHeader("Authorization")
		if !stringx.IsHasPrefix(authorization, "Bearer ") {
			ClientUnauthorizedResponse(ginCtx, nil, "Authorization header does not start with 'Bearer '")
			return
		}

		token := stringx.TrimPrefix(authorization, "Bearer ")
		if stringx.IsEmptyStr(token) {
			ClientUnauthorizedResponse(ginCtx, nil, "Token is empty")
			return
		}

		conn, err := h.grpcPools.MicroAuthGrpcConnPool.GetConnFromPool()
		if errorx.IsErrorExist(err) {
			ServerInternalErrorResponse(ginCtx, nil, "h.grpcPools.MicroAuthGrpcConnPool.GetConnFromPool err")
			return
		}
		grpcResp, err := micro_auth.NewAuthServiceClient(conn).ValidateToken(
			ginCtx.Request.Context(),
			&micro_auth.ValidateTokenRequest{
				Token: token,
			},
		)
		if errorx.IsErrorExist(err) {
			ClientUnauthorizedResponse(ginCtx, nil, fmt.Sprintf("ValidateToken err: %s", err.Error()))
			return
		}
		if !grpcResp.Valid {
			ClientUnauthorizedResponse(ginCtx, nil, "Token is invalid")
			return
		}

		ginCtx.Request = ginCtx.Request.WithContext(
			context.WithValue(
				ginCtx.Request.Context(),
				TokenInfoCtxKey,
				TokenInfo{
					ID:    grpcResp.ID,
					Email: grpcResp.Email,
				},
			),
		)

		// 處理請求
		ginCtx.Next()
	}
}
