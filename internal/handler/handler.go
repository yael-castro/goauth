package handler

import (
	"github.com/yael-castro/goauth/internal/business"
	"net/http"
	"net/url"
)

// OAuthError takes a *url.URL to set OAuth errors in their query parameters
func OAuthError(uri *url.URL, err error, description string) {
	q := url.Values{}

	q.Set("error", err.Error())
	q.Set("error_description", description)

	uri.RawQuery = q.Encode()
}

// NewServeMux builds a http.ServeMux based on a business.CodeGrant
// and is returned as http.Handler
func NewServeMux(grant business.CodeGrant) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/go-auth/v1/authorization", NewAuthorizationHandler(grant))
	mux.HandleFunc("/go-auth/v1/token", NewTokenHandler(grant))

	return mux
}
