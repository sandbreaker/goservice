package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/sandbreaker/goservice/config"
	"github.com/sandbreaker/goservice/log"
	"github.com/sandbreaker/goservice/service"
)

func init() {
}

func main() {
	flag.Parse()
	defer sentry.Flush(2 * time.Second)

	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	api, err := service.NewServiceAPI()
	if err != nil {
		log.Fatal(err)
	}
	defer api.Close()

	role := config.StaticConfig.GetRole()
	port := config.StaticConfig.GetInt("goservice.api.port", 8088)

	// just use :<port> if global listen
	var listenHost string
	if config.StaticConfig.GetEnv() == config.EnvDev {
		listenHost = fmt.Sprintf("localhost:%d", port)
	} else {
		listenHost = fmt.Sprintf(":%d", port)
	}

	// init sentry
	setSentry(config.StaticConfig.GetString("sentry.dsn."+role, ""), true)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      api.Router,
		Addr:         listenHost,
	}

	log.Info("Starting service on host=", listenHost, ", role=", role, ", env=", config.StaticConfig.GetEnv())

	srv.ListenAndServe()
}

func setSentry(dsn string, debug bool) {
	// Either set your DSN here or set the SENTRY_DSN environment variable.
	err := sentry.Init(sentry.ClientOptions{
		Dsn:   dsn,
		Debug: debug,
	})
	if err != nil {
		log.Fatal("Sentry fail to init: ", err)
	}
}
