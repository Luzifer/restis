package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		DisableCORS     bool   `flag:"disable-cors" default:"false" description:"Disable setting CORS headers for all requests"`
		Listen          string `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel        string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		RedisConnString string `flag:"redis-conn-string" default:"redis://localhost:6379/0" description:"Connection string for redis"`
		RedisKeyPrefix  string `flag:"redis-key-prefix" default:"" description:"Prefix to prepend to keys (will be prepended without delimiter!)"`
		VersionAndExit  bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func initApp() error {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		return errors.Wrap(err, "parsing cli options")
	}

	if cfg.VersionAndExit {
		fmt.Printf("git-changerelease %s\n", version)
		os.Exit(0)
	}

	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return errors.Wrap(err, "parsing log-level")
	}
	logrus.SetLevel(l)

	return nil
}

func main() {
	var err error
	if err = initApp(); err != nil {
		logrus.WithError(err).Fatal("initializing app")
	}

	redisOpts, err := redis.ParseURL(cfg.RedisConnString)
	if err != nil {
		logrus.WithError(err).Fatal("parsing redis connection string")
	}

	var (
		client = redis.NewClient(redisOpts)
		router = mux.NewRouter()
	)

	if !cfg.DisableCORS {
		router.Use(corsMiddleware)
	}

	router.MethodNotAllowedHandler = corsMiddleware(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			// Most likely JS client asking for CORS headers
			res.WriteHeader(http.StatusNoContent)
			return
		}

		res.WriteHeader(http.StatusMethodNotAllowed)
	}))

	router.Methods(http.MethodDelete).HandlerFunc(handlerDelete(client))
	router.Methods(http.MethodGet).HandlerFunc(handlerGet(client))
	router.Methods(http.MethodPut).HandlerFunc(handlerPut(client))

	server := &http.Server{
		Addr:              cfg.Listen,
		Handler:           router,
		ReadHeaderTimeout: time.Second,
	}

	logrus.WithFields(logrus.Fields{
		"addr":    cfg.Listen,
		"version": version,
	}).Info("starting HTTP server")

	if err := server.ListenAndServe(); err != nil {
		logrus.WithError(err).Fatal("HTTP server quit unexpectedly")
	}
}
