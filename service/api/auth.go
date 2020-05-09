package api

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"github.com/alexedwards/scs"
	"github.com/gorilla/securecookie"
	"github.com/palantir/go-githubapp/oauth2"
	"github.com/rs/zerolog"
	oauth22 "golang.org/x/oauth2"
	"net/http"
	"time"
)

const (
	AuthRoute       = "/api/auth/github"
	TokenCookieName = "token"
)

type OAuthManager struct {
	Config       Config
	SecureCookie *securecookie.SecureCookie
}

func (o OAuthManager) GetHandler() http.Handler {
	sessionManager := scs.NewCookieManager(o.Config.AppConfig.CookieSecret)
	store := oauth2.SessionStateStore{Sessions: sessionManager}

	return oauth2.NewHandler(
		GetOAuthConfig(o.Config),
		oauth2.ForceTLS(!o.Config.AppConfig.Debug),
		oauth2.WithStore(&store),
		oauth2.OnLogin(func(w http.ResponseWriter, r *http.Request, login *oauth2.Login) {
			o.saveTokenToCookie(login.Token, &w)

			http.Redirect(w, r, o.Config.AppConfig.DashboardUrl, http.StatusFound)
		}),
	)
}

func (o OAuthManager) saveTokenToCookie(token *oauth22.Token, w *http.ResponseWriter) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(token)
	if err != nil {
		panic("failed to encode token")
	}

	tokenString := base64.StdEncoding.EncodeToString(buffer.Bytes())
	encoded, err := o.SecureCookie.Encode(TokenCookieName, tokenString)
	if err != nil {
		panic("failed to encode token base 64")
	}

	cookie := &http.Cookie{
		Name:  TokenCookieName,
		Value: encoded,
		Path:  "/",
	}
	http.SetCookie(*w, cookie)
}

func GetOAuthConfig(config Config) *oauth22.Config {
	return oauth2.GetConfig(config.Github, []string{"user:email"})
}

func GetTokenFromRequest(r *http.Request, secureCookie *securecookie.SecureCookie) (*oauth22.Token, error) {
	tokenCookieBase64, err := r.Cookie(TokenCookieName)
	if err != nil {
		return nil, errors.New("")
	}

	decodedTokenCookie := ""
	err = secureCookie.Decode(TokenCookieName, tokenCookieBase64.Value, &decodedTokenCookie)
	if err != nil {
		return nil, errors.New("failed to decode secured cookie")
	}

	by, err := base64.StdEncoding.DecodeString(decodedTokenCookie)
	if err != nil {
		return nil, errors.New("failed to decode cookie base 64")
	}
	token := oauth22.Token{}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&token)
	if err != nil {
		return nil, errors.New("failed to decode token")
	}
	return &token, err
}

type IsAuthenticatedHandler struct {
	SecureCookie *securecookie.SecureCookie
}

func (h IsAuthenticatedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())

	result := IsAuthenticated{IsAuthenticated: false}
	token, err := GetTokenFromRequest(r, h.SecureCookie)
	if err != nil {
		if len(err.Error()) > 0 {
			logger.Error().Err(err)
		}
	} else {
		result.IsAuthenticated = token.Expiry.Before(time.Now())
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		panic(err)
	}
}

type IsAuthenticated struct {
	IsAuthenticated bool `json:"is_authenticated"`
}

type SignOutHandler struct {
	Config Config
}

func (h SignOutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{
		Name:    TokenCookieName,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, c)
	http.Redirect(w, r, h.Config.AppConfig.DashboardUrl, http.StatusFound)
}
