package handler

import (
	"context"
	"errors"
	"net/http"
	"tools/oauth"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// Login by google oauth2
func Login(o *oauth2.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, o.AuthCodeURL(oauth.State))
		c.JSON(http.StatusOK, gin.H{
			"message": "Login",
		})
	}
}

func GoogleCallback(o *oauth2.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.Query(oauth.State)
		if s != oauth.State {
			c.AbortWithError(http.StatusUnauthorized, errors.New("state is not valid"))
			return
		}

		code := c.Query(oauth.Code)
		token, err := o.Exchange(context.Background(), code)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": token.AccessToken,
		})
	}
}
