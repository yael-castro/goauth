// Package handler handles all http requests made to the app (is the presentation layer)
package handler

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"path"

	"github.com/yael-castro/godi/internal/model"
)

// NotFound handles not found requests
func NotFound(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusNotFound, model.Map{"message": "not found"})
}

// MethodNotAllowed handles invalid requests using an illegal method
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write(nil)
}

// Healthcheck handles requests made to check the status server
func Healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		MethodNotAllowed(w, r)
		return
	}

	JSON(w, http.StatusOK, model.Map{"message": "ok"})
}

// JSON function used to send serialized data in json more easier
func JSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func Bind(r *http.Request, i interface{}) error {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	switch contentType {
	case "application/json":
		return json.NewDecoder(r.Body).Decode(i)
	}

	return fmt.Errorf("mime type not supported '%s'", contentType)
}

// New constructs an empty Handler
func New() *Handler {
	return &Handler{}
}

// Handler main handler used in the ListeAndServe
type Handler struct {
	Authenticator http.Handler
	Authorizator  http.Handler
	Revoker       http.Handler
}

// ServeHTTP decides which http.HandlerFunc use based on the http method
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := path.Join(r.URL.Path, "/")

	switch p {
	case "/goauth/v1/authenticate":
		h.Authenticator.ServeHTTP(w, r)

	case "/goauth/v1/authorizate":
		h.Authorizator.ServeHTTP(w, r)

	case "/goauth/v1/revoke":
		h.Revoker.ServeHTTP(w, r)

	case "/goauth/v1/healthcheck":
		Healthcheck(w, r)

	default:
		NotFound(w, r)
	}
}
