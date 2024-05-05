package handler

import (
	"net/http"
	"tools/oauth"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Login by google oauth2
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := oauth.CreateGoogleOauth(
			viper.GetString("oauth2.google.client_id"),
			viper.GetString("oauth2.google.client_secret"),
			viper.GetString("oauth2.google.redirect_url"),
			viper.GetStringSlice("oauth2.google.scopes"),
		)

		c.Redirect(http.StatusSeeOther, url)

		//s := c.Query("state")
		//if s != "state_parameter_passthrough_value" {
		//	c.AbortWithError(http.StatusUnauthorized,, "error message")
		//	return
		//}
		//
		//code := c.Query("code")
		//token, err := config.Exchange(oauth2.NoContext, code)

		c.JSON(http.StatusOK, gin.H{
			"message": "Login",
		})
	}
}
