package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func CreateGoogleOauth(id, secret, redirectURL string, scopes []string) string {
	c := oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}

	return c.AuthCodeURL("state")
}
