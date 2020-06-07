//  Copyright © 2020 Sang Chi. All rights reserved.

package service

import (
	"net/http"

	"github.com/sandbreaker/goservice/modeldefault"
)

type Context struct {
	*ServiceAPI
	errorTags map[string]string

	// default
	UserDefault *modeldefault.User
}

func DefaultAppContext(api *ServiceAPI, r *http.Request, user *modeldefault.User) *Context {
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
