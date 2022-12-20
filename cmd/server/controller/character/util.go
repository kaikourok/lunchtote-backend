package character

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/library/secure"
	"golang.org/x/oauth2"
)

func (u *CharacterController) oauthFailedResponse(c *gin.Context, mode string, state int) {
	config := u.registry.GetConfig()

	if mode != "signin" && mode != "register" {
		c.String(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest)+" Bad Request.")
		return
	}

	if mode == "signin" {
		c.Redirect(http.StatusFound, config.GetString("oauth.signin-url")+"?result="+strconv.Itoa(state))
	} else {
		c.Redirect(http.StatusFound, config.GetString("oauth.registered-uri")+"?result="+strconv.Itoa(state))
	}
}

func (u *CharacterController) getSessionDefaultConfig() sessions.Options {
	config := u.registry.GetConfig()

	return sessions.Options{
		MaxAge: config.GetInt("secure.max-age"),
		Path:   config.GetString("general.api-path"),
	}
}

func (u *CharacterController) getGoogleAuthConfig() *oauth2.Config {
	config := u.registry.GetConfig()

	return &oauth2.Config{
		ClientID:     config.GetString("oauth.google.client-id"),
		ClientSecret: config.GetString("oauth.google.client-secret"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.GetString("oauth.google.auth-url"),
			TokenURL: config.GetString("oauth.google.token-url"),
		},
		RedirectURL: config.GetString("general.host") + config.GetString("oauth.google.redirect-uri"),
		Scopes:      config.GetStringSlice("oauth.google.scopes"),
	}
}

func (u *CharacterController) getTwitterAuthConfig() *oauth2.Config {
	config := u.registry.GetConfig()

	return &oauth2.Config{
		ClientID:     config.GetString("oauth.twitter.client-id"),
		ClientSecret: config.GetString("oauth.twitter.client-secret"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.GetString("oauth.twitter.auth-url"),
			TokenURL: config.GetString("oauth.twitter.token-url"),
		},
		RedirectURL: config.GetString("general.host") + config.GetString("oauth.twitter.redirect-uri"),
		Scopes:      config.GetStringSlice("oauth.twitter.scopes"),
	}
}

func (u *CharacterController) generateCodeVerifier() string {
	config := u.registry.GetConfig()
	return secure.GenerateSecureRandomHex(config.GetInt("oauth.verifier-length"))
}

func (u *CharacterController) convertCodeChallenge(verifier string) string {
	s := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(s[:])
}
