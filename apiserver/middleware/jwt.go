package middleware

import (
	"context"
	"net/http"

	"definition/api_response"

	"github.com/gin-gonic/gin"

	"tools/encryption"
	"tools/logger"
	"tools/stringx"
)

const TokenInfoCtxKey = "TokenInfo"

type TokenInfo struct {
	InvestorID string
}

// TODO 加入 expiretime and 將Token 存入 redis, 我要做到logout, 清除所有這個user的token
func JWT(
	encryption encryption.IEncryption[
		encryption.JWTRequirements,
		encryption.JWTMapClaims, string,
		string, TokenInfo,
	],
) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("JWTMiddleware Start")
		// Receive request Header Token
		authorization := c.GetHeader("Authorization")

		if !stringx.CheckHasPrefix(authorization, "Bearer ") {
			logger.Error("Authorization header does not start with 'Bearer '")
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				api_response.SetAPIFormatResponse("Invalid token", api_response.ClientUnauthorized, nil),
			)
			return
		}

		// Remove Bearer prefix
		token := stringx.TrimPrefix(authorization, "Bearer ")

		if stringx.CheckEmptyStr(token) {
			logger.Error("JWTMiddleware Token is empty")
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				api_response.SetAPIFormatResponse("Invalid token", api_response.ClientUnauthorized, nil),
			)
			return
		}

		decrypted, err := encryption.Decrypt(token)
		if err != nil {
			logger.Error("JWTMiddleware Decrypt error: %v", err)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				api_response.SetAPIFormatResponse("Invalid token", api_response.ClientUnauthorized, nil),
			)
			return
		}

		// Add decrypted token to context
		c.Request = c.Request.WithContext(
			context.WithValue(
				c.Request.Context(),
				TokenInfoCtxKey,
				decrypted,
			),
		)

		// 處理請求
		c.Next()
	}
}

// TODO 準備實作, 如何利用Operation name 穿過JWT 驗證, 現在Register 和 Login 都是不需要JWT 驗證的
/*
func JWTMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Assume the operation name is passed in the query string for simplicity
    operationName := r.URL.Query().Get("operationName")

    // Skip authentication for login and register
    if operationName == "Login" || operationName == "Register" {
      next.ServeHTTP(w, r)
      return
    }

    token := extractTokenFromHeader(r)
    if token == "" {
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
    }

    // Verify token logic here
    if !jwt.VerifyToken(token) {
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
    }

    ctx := context.WithValue(r.Context(), "userID", jwt.GetUserIDFromToken(token))
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}
*/
