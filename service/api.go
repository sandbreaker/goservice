//  Copyright © 2020 Sang Chi. All rights reserved.

package service

import (
	"errors"
	"net/http"
	"net/http/pprof"

	"github.com/sandbreaker/goservice/modeldefault"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/palantir/stacktrace"
)

const (
	DefaultMaxFormMemory = 32 << 20 // 32 MB
)

var decoder = schema.NewDecoder()

type handleFunc func(w http.ResponseWriter, r *http.Request, ctx *Context)

type authLevel int

const (
	AuthLevelAny authLevel = iota
	AuthLevelLoggedIn
	AuthLevelAdmin
)

type productType int

const (
	ProdTypeDefault productType = iota
)

type ServiceAPI struct {
	Name    string
	Version string
	Router  *mux.Router
}

func NewServiceAPI() (*ServiceAPI, error) {
	var err error
	api := &ServiceAPI{}

	api.initRouter()

	return api, err
}

func (*ServiceAPI) Close() {
}

func (api *ServiceAPI) getProdTypeName(prodType productType) string {
	switch prodType {
	case ProdTypeDefault:
		return "default"
	default:
		return "unknown"
	}
}

func (api *ServiceAPI) middleware(h handleFunc, apiName string, prodType productType, level authLevel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(DefaultMaxFormMemory)
		// var prodName = api.getProdTypeName(prodType)
		var err error

		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = stacktrace.Propagate(errors.New(t), "")
				case error:
					err = stacktrace.Propagate(t, "")
				default:
					err = errors.New("What? No valid error found")
				}

				// metric.StaticMetric.Increment(fmt.Sprintf("%s.%s.%s", prodName, metric.MetricAPIErr, apiName), 1, nil)
				// util.LogAndCaptureError(make(map[string]string), stacktrace.Propagate(err, "PANIC PANIC PANIC stack: %s", debug.Stack()))
				// service.InternalServerError(w, err)
				return
			}
		}()

		// preflight check
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			headers := w.Header()
			headers.Add("Vary", "Origin")
			headers.Add("Vary", "Access-Control-Request-Method")
			headers.Add("Vary", "Access-Control-Request-Headers")

			headers.Set("Access-Control-Allow-Origin", "*")
			headers.Set("Access-Control-Allow-Methods", "POST,PUT,PATCH,GET,DELETE,OPTIONS")
			headers.Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,x-access-token")

			HandleOk(w, nil)
			return
		}

		// for access control from all domains
		// ONLY do this for testing for js applications not on same domain
		if r.Method == http.MethodGet || r.Method == http.MethodPost {
			headers := w.Header()
			headers.Set("Access-Control-Allow-Origin", "*")
		}

		switch prodType {
		case ProdTypeDefault:
			user := &modeldefault.User{}
			h(w, r, DefaultAppContext(api, r, user))

		}
	}
}

func (api *ServiceAPI) initRouter() {
	api.Router = mux.NewRouter()

	// sub-route
	debugRoute := api.Router.PathPrefix("/pprof").Subrouter()
	utilRoute := api.Router.PathPrefix("/util").Subrouter()
	// apiRoute := api.Router.PathPrefix("/api").Subrouter()

	// pprof api, need security
	debugRoute.Handle("/pprof/", http.HandlerFunc(pprof.Index))
	debugRoute.Handle("/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	debugRoute.Handle("/pprof/profile", http.HandlerFunc(pprof.Profile))
	debugRoute.Handle("/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	debugRoute.Handle("/pprof/trace", http.HandlerFunc(pprof.Trace))
	debugRoute.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	debugRoute.Handle("/pprof/heap", pprof.Handler("heap"))
	debugRoute.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	debugRoute.Handle("/pprof/block", pprof.Handler("block"))

	// utility routes
	utilRoute.Handle("/version", api.middleware(GetVersion, "getVersion", ProdTypeDefault, AuthLevelAny)).Methods("GET")

}
