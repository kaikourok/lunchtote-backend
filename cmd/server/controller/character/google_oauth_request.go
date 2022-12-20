package character

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/library/secure"
	"golang.org/x/oauth2"
)

func (u *CharacterController) GoogleOauthRequest(c *gin.Context) {
	session := sessions.Default(c)
	config := u.registry.GetConfig()

	mode := c.Query("mode")
	if mode != "signin" && mode != "register" {
		u.oauthFailedResponse(c, mode, 0)
		return
	}

	if mode == "register" {
		if session.Get("cid") == nil {
			u.oauthFailedResponse(c, "signin", http.StatusInternalServerError)
			return
		}
	}

	authConfig := u.getGoogleAuthConfig()
	state := secure.GenerateSecureRandomHex(config.GetInt("oauth.state-length"))

	codeVerifier := u.generateCodeVerifier()
	codeChallenge := u.convertCodeChallenge(codeVerifier)

	session.Set("google-state", state)
	session.Set("google-oauth-mode", mode)
	session.Set("google-code-verifier", codeVerifier)
	session.Options(u.getSessionDefaultConfig())
	session.Save()

	challengeOption := oauth2.SetAuthURLParam("code_challenge", codeChallenge)
	challengeMethodOption := oauth2.SetAuthURLParam("code_challenge_method", "S256")

	authCodeURL := authConfig.AuthCodeURL(state, challengeOption, challengeMethodOption)

	c.Redirect(http.StatusFound, authCodeURL)
}
