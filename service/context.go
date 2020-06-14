package service

import (
	"net/http"

	"github.com/sandbreaker/goservice/datamodel/mdefault"
)

type Context struct {
	*ServiceAPI
	errorTags map[string]string

	// default
	UserDefault *mdefault.User
}

func DefaultAppContext(api *ServiceAPI, r *http.Request, user *mdefault.User) *Context {
	ctx := &Context{
		ServiceAPI:  api,
		UserDefault: user,
		errorTags:   make(map[string]string),
	}

	if user != nil {
		ctx.errorTags["uid"] = user.UID
	}

	return ctx
}
