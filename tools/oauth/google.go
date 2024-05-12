package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	Code  = "code"
	State = "state"
)

func NewGoogleOauth(id, secret, redirectURL string, scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
}
