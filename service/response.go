//  Copyright © 2020 Sang Chi. All rights reserved.

package service

import (
	"net/http"

	"github.com/sandbreaker/goservice/log"

	"github.com/bitly/go-simplejson"
)

func SetContentTypeJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func SetContentTypeHtml(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func SetAccessControlOrigin(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func writeError(w http.ResponseWriter, httpCode int, err error) {
	SetContentTypeJson(w)
	w.WriteHeader(httpCode)
	w.Write(errorResponse(httpCode, err))
}

func HandleOk(w http.ResponseWriter, payload []byte) {
	SetContentTypeJson(w)
	w.WriteHeader(http.StatusOK)

	if payload != nil {
		w.Write(payload)
	} else {
		w.Write(successResponse())
	}
}

func InternalServerError(w http.ResponseWriter, err error) {
	writeError(w, http.StatusInternalServerError, err)
}

func UnauthorizedAccessError(w http.ResponseWriter, err error) {
	writeError(w, http.StatusUnauthorized, err)
}

func NotFoundError(w http.ResponseWriter, err error) {
	writeError(w, http.StatusNotFound, err)
}

func successResponse() []byte {
	json := simplejson.New()
	json.Set("code", 200)
	json.Set("result", "success")

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Error(err)
		return nil
	}

	return payload
}

func errorResponse(code int, err error) []byte {
	json := simplejson.New()
	json.Set("code", code)

	if err != nil {
		json.Set("error", err.Error())
		// TODO, arg, this needs to work like success response to be consistent, write twice for now
		json.Set("result", err.Error())
	}

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Error(err)
		return nil
	}

	return payload
}
