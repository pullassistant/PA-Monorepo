package api

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/gorilla/securecookie"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
	"net/http"
)

type clientKey struct{}
type appClientKey struct{}

func clientFromCtx(context context.Context) github.Client {
	if t, ok := context.Value(clientKey{}).(github.Client); ok {
		return t
	}

	panic("No client in context")
}

func appClientFromCtx(context context.Context) github.Client {
	if t, ok := context.Value(appClientKey{}).(github.Client); ok {
		return t
	}

	panic("No app client in context")
}

func EnsureToken(h http.Handler, secureCookie *securecookie.SecureCookie, config Config, clientCreator githubapp.ClientCreator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.Ctx(r.Context())

		token, err := GetTokenFromRequest(r, secureCookie)
		if err != nil {
			if len(err.Error()) > 0 {
				logger.Error().Err(err)
			}
			http.Redirect(w, r, AuthRoute, http.StatusFound)
			return
		}

		client := GetOAuthConfig(config).Client(r.Context(), token)
		githubClient := github.NewClient(client)

		appClient, err := clientCreator.NewAppClient()
		if err != nil {
			panic(err)
		}

		newCtx := context.WithValue(r.Context(), clientKey{}, *githubClient)
		newCtx = context.WithValue(newCtx, appClientKey{}, *appClient)

		h.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func EnsureAllowOrigin(h http.Handler, config Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.AppConfig.Debug {
			w.Header().Set("Access-Control-Allow-Origin", config.AppConfig.DashboardUrl)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		h.ServeHTTP(w, r)
	})
}

func EnsureTLS(h http.Handler, config Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !config.AppConfig.Debug {
			forwardedProtocol := r.Header.Get("x-forwarded-proto")
			if len(forwardedProtocol) == 0 {
				h.ServeHTTP(w, r)
			} else if forwardedProtocol == "https" {
				h.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, config.AppConfig.DashboardUrl, http.StatusFound)
			}
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

//func getOrganizationClient(cc githubapp.ClientCreator, orgName string, context context.Context) (*github.Client, error)r {
//	appClient, err := cc.NewAppClient()
//	if err != nil {
//		return nil, err
//	}
//
//	installations := githubapp.NewInstallationsService(appClient)
//	installation, err := installations.GetByOwner(context, orgName)
//	if err != nil {
//		panic(err)
//	}
//
//	return cc.NewInstallationClient(installation.ID)
//}
