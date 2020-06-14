//  Copyright © 2020 Sang Chi. All rights reserved.

package service

import "net/http"

func GetHealth(w http.ResponseWriter, r *http.Request, ctx *Context) {
	HandleOk(w, nil)
}

func GetVersion(w http.ResponseWriter, r *http.Request, ctx *Context) {
	HandleOk(w, nil)
}
