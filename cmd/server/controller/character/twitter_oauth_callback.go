package character

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kaikourok/lunchtote-backend/library/secure"
	"golang.org/x/oauth2"
)

func (u *CharacterController) TwitterOauthCallback(c *gin.Context) {
	session := sessions.Default(c)
	config := u.registry.GetConfig()

	stateInterface := session.Get("twitter-state")
	modeInterface := session.Get("twitter-oauth-mode")
	verifierInterface := session.Get("twitter-code-verifier")

	if stateInterface == nil || modeInterface == nil || verifierInterface == nil {
		u.oauthFailedResponse(c, "", 0)
		return
	}
	state := stateInterface.(string)
	mode := modeInterface.(string)
	codeVerifier := verifierInterface.(string)

	respondState := c.Query("state")
	if respondState == "" {
		u.oauthFailedResponse(c, mode, http.StatusUnauthorized)
		return
	}

	code := c.Query("code")
	if code == "" {
		u.oauthFailedResponse(c, mode, http.StatusUnauthorized)
		return
	}

	if state != respondState {
		u.oauthFailedResponse(c, mode, http.StatusUnauthorized)
		return
	}

	authConfig := u.getTwitterAuthConfig()

	verifierOption := oauth2.SetAuthURLParam("code_verifier", codeVerifier)
	token, err := authConfig.Exchange(context.Background(), code, verifierOption)

	if err != nil {
		u.oauthFailedResponse(c, mode, http.StatusInternalServerError)
		return
	}

	request, err := http.NewRequest("GET", config.GetString("oauth.twitter.user-data-url"), nil)
	request.Header.Set("Authorization", "Bearer "+token.AccessToken)

	if err != nil {
		u.oauthFailedResponse(c, mode, http.StatusInternalServerError)
		return
	}

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		u.oauthFailedResponse(c, mode, http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	if err != nil {
		u.oauthFailedResponse(c, mode, http.StatusInternalServerError)
		return
	}

	type TwitterUserDataResponse struct {
		Data struct {
			Id       string `json:"id"`
			Name     string `json:"name"`
			UserName string `json:"username"`
		} `json:"data"`
	}

	var userData TwitterUserDataResponse
	err = json.Unmarshal(body, &userData)
	if err != nil {
		u.oauthFailedResponse(c, mode, http.StatusInternalServerError)
		return
	}

	if mode == "signin" {
		characterId, notificationToken, err := u.usecase.RetrieveCredentialsByTwitter(userData.Data.Id)
		if err != nil {
			u.oauthFailedResponse(c, mode, http.StatusUnauthorized)
			return
		}

		csrfToken := secure.GenerateSecureRandomHex(config.GetInt("secure.token-length"))

		session.Clear()
		session.Set("cid", characterId)
		session.Set("csrf-token", csrfToken)
		session.Set("notification-token", notificationToken)
		session.Options(u.getSessionDefaultConfig())
		session.Save()

		c.Redirect(http.StatusFound, config.GetString("general.host")+config.GetString("oauth.signed-in-uri"))
		return
	} else if mode == "register" {
		session.Delete("twitter-state")
		session.Delete("twitter-oauth-mode")
		session.Delete("twitter-code-verifier")
		session.Options(u.getSessionDefaultConfig())
		session.Save()

		err = u.usecase.RegisterTwitterData(session.Get("cid").(int), userData.Data.Id)
		if err != nil {
			u.oauthFailedResponse(c, mode, http.StatusInternalServerError)
			return
		}

		c.Redirect(http.StatusFound, config.GetString("general.host")+config.GetString("oauth.registered-uri"))
		return
	} else {
		u.oauthFailedResponse(c, "", 0)
		return
	}
}
