//  Copyright © 2020 Sang Chi. All rights reserved.

package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/sandbreaker/goservice/service"
)

func init() {
}

func main() {
	flag.Parse()

	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	api, err := service.NewServiceAPI()
	if err != nil {
		log.Fatal(err)
	}
	defer api.Close()

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      api.Router,
		Addr:         "localhost:8080",
	}

	srv.ListenAndServe()

}
