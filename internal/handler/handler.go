// Package handler handles all http requests made to the app (is the presentation layer)
package handler

import (
	"net/http"
)

// MethodNotAllowed handles invalid requests using an illegal method
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write(nil)
}

// HealthCheck handles requests made to check the status server
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		MethodNotAllowed(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

// NotFound handles request made to not found source
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(nil)
}
