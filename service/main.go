package main

import (
	"bytes"
	"github.com/domsu/pullassistant/api"
	"github.com/gorilla/securecookie"
	"github.com/palantir/go-baseapp/baseapp"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
	"goji.io/pat"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const agentVersionName = "pull-assistant/1.0.0"

func main() {
	configFile := "config.yml"
	if v, ok := os.LookupEnv("PULL_ASSISTANT_ALPHA"); ok {
		if v == "true" {
			configFile = "config_alpha.yml"
		}
	}

	config, err := api.ReadConfig(os.Getenv("KO_DATA_PATH") + "/config/" + configFile)
	if err != nil {
		panic(err)
	}

	tracker := newTracker(config.AppConfig.TrackingId, config.AppConfig.AppName)
	jobChan := make(chan PRWorkerData)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	startPRWorker(jobChan, logger, tracker, config)
	startServer(jobChan, logger, config)
}

func startPRWorker(jobChan <-chan PRWorkerData, logger zerolog.Logger, tracker *Tracker, config *api.Config) {
	worker := PRWorker{jobChan: jobChan, logger: logger, tracker: tracker, config: config}
	worker.start()
}

func startServer(jobChan chan<- PRWorkerData, logger zerolog.Logger, config *api.Config) {
	server, err := baseapp.NewServer(
		config.Server,
		baseapp.DefaultParams(logger, "")...,
	)
	if err != nil {
		panic(err)
	}

	var middleware githubapp.ClientOption = nil
	if config.AppConfig.Debug {
		middleware = githubapp.WithClientMiddleware(
			githubapp.ClientMetrics(server.Registry()),
			HttpLogger(zerolog.NoLevel),
		)
	} else {
		middleware = githubapp.WithClientMiddleware(
			githubapp.ClientMetrics(server.Registry()),
			githubapp.ClientLogging(zerolog.NoLevel),
		)
	}

	clientCreator, err := githubapp.NewDefaultCachingClientCreator(
		config.Github,
		githubapp.WithClientUserAgent(agentVersionName),
		middleware,
	)
	if err != nil {
		panic(err)
	}

	var secureCookie = securecookie.New([]byte(config.AppConfig.CookieSecret), []byte(config.AppConfig.CookieSecret))

	prCommentHandler := &PRHandler{
		ClientCreator: clientCreator,
		channel:       jobChan,
	}
	webHookHandler := githubapp.NewDefaultEventDispatcher(config.Github, prCommentHandler)

	server.Mux().Handle(pat.Post(githubapp.DefaultWebhookRoute), webHookHandler)

	handler := api.OAuthManager{*config, secureCookie}
	server.Mux().Handle(pat.Get(api.AuthRoute), handler.GetHandler())
	server.Mux().Handle(pat.Get("/api/installations"), api.EnsureAllowOrigin(api.EnsureToken(api.InstallationsHandler{*config}, secureCookie, *config, clientCreator), *config))
	server.Mux().Handle(pat.Get("/api/repositories/:installationId"), api.EnsureAllowOrigin(api.EnsureToken(api.RepositoryHandler{}, secureCookie, *config, clientCreator), *config))
	server.Mux().Handle(pat.Get("/api/auth/isAuthenticated"), api.EnsureAllowOrigin(api.IsAuthenticatedHandler{secureCookie}, *config))
	server.Mux().Handle(pat.Get("/api/auth/signOut"), api.SignOutHandler{*config})
	server.Mux().Handle(pat.Get("/api/dashboardConfig"), api.EnsureAllowOrigin(api.DashboardConfigHandler{*config}, *config))

	server.Mux().Handle(pat.Get("/*"), api.EnsureTLS(http.FileServer(http.Dir(os.Getenv("KO_DATA_PATH")+"/dashboard/")), *config))

	_ = server.Start()
}

func HttpLogger(lvl zerolog.Level) githubapp.ClientMiddleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			start := time.Now()
			res, err := next.RoundTrip(r)
			elapsed := time.Now().Sub(start)

			log := zerolog.Ctx(r.Context())
			if res != nil {
				buf := new(bytes.Buffer)
				content := ""
				if b, err := ioutil.ReadAll(buf); err == nil {
					content = string(b)
				}

				log.WithLevel(lvl).
					Str("method", r.Method).
					Str("path", r.URL.String()).
					Int("status", res.StatusCode).
					Str("content", content).
					Int64("size", res.ContentLength).
					Dur("elapsed", elapsed).
					Msg("github_request")
			} else {
				log.WithLevel(lvl).
					Str("method", r.Method).
					Str("path", r.URL.String()).
					Int("status", -1).
					Int64("size", -1).
					Dur("elapsed", elapsed).
					Msg("github_request")
			}

			return res, err
		})
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}
