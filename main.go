package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	chilogger "github.com/igknot/chi-zap-ecs-logger"
	"github.com/ory/graceful"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")
	creds    = map[string]string{username: password}
)

func main() {
	// create new Client
	c, err := NewClient(cfenv.IsRunningOnCF(), os.Stdout, zapcore.InfoLevel)
	if err != nil {
		log.Printf("could not create client: %v", err)
		os.Exit(-1)
	}
	defer c.Logger.Sync()

	r := chi.NewRouter()
	r.Use(middleware.CleanPath)
	r.Use(middleware.Heartbeat("/healthz"))
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.StripSlashes)
	r.Use(chilogger.NewZapMiddleware("go-cfserver", c.Logger))

	r.Route("/", func(r chi.Router) {
		r.Handle("/actuator/prometheus", promhttp.Handler())

		// basic auth section
		r.With(middleware.BasicAuth("secret", creds)).Get("/env", c.GetEnv)
	})

	// gracefull shutdown
	srv := graceful.WithDefaults(&http.Server{
		Addr:    ":8080",
		Handler: r,
	})

	c.Logger.Info("starting the server")
	if err := graceful.Graceful(srv.ListenAndServe, srv.Shutdown); err != nil {
		c.Logger.Error("failed to gracefully shutdown", zap.Error(err))
		os.Exit(-1)
	}
	c.Logger.Info("server was shutdown gracefully")
}
