package service

import (
	"net/http"

	"github.com/bitly/go-simplejson"

	"github.com/sandbreaker/goservice/log"
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

func HandleOk(w http.ResponseWriter, payload []byte) {
	SetContentTypeJson(w)
	w.WriteHeader(http.StatusOK)

	if payload != nil {
		w.Write(payload)
	} else {
		w.Write(successResponse())
	}
}

func EroorInternalServer(w http.ResponseWriter, err error) {
	writeError(w, http.StatusInternalServerError, err)
}

func ErrorNotFound(w http.ResponseWriter, err error) {
	writeError(w, http.StatusNotFound, err)
}

func ErrorUnauthorizedAccess(w http.ResponseWriter, err error) {
	writeError(w, http.StatusUnauthorized, err)
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
		// write twice for now to work like success response for consistency
		json.Set("result", err.Error())
	}

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Error(err)
		return nil
	}

	return payload
}

func writeError(w http.ResponseWriter, httpCode int, err error) {
	SetContentTypeJson(w)
	w.WriteHeader(httpCode)
	w.Write(errorResponse(httpCode, err))
}
