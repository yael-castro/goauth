package business

import (
	"fmt"
	"net/url"

	"github.com/yael-castro/godi/internal/model"
)

// Authorizer handles authorization requests for any flow of the OAuth 2.0 protocol
type Authorizer interface {
	// Authorize receives an authorization request and returns the redirect uri
	Authorize(model.Authorization) *url.URL
}

// OAuth hash map to manage the different flows of the OAuth 2.0 protocol
type OAuth map[string]Authorizer

// Authorize selects an Authorizer based on the flow specified to be executed
func (o OAuth) Authorize(auth model.Authorization) *url.URL {
	uri := auth.RedirectURL

	if auth.ResponseType == "" {
		return OAuthError(uri, model.UnsupportedResponseType, "response type is not specified in the requests")
	}

	flow, ok := o[auth.ResponseType]
	if !ok {
		description := fmt.Sprintf(`response_type "%s" is not supported by the server`, auth.ResponseType)
		return OAuthError(uri, model.UnsupportedResponseType, description)
	}

	return flow.Authorize(auth)
}

// OAuthError takes a *url.URL to set OAuth errors in their query parameters
func OAuthError(uri *url.URL, err model.OAuthError, description string) *url.URL {
	q := url.Values{}

	q.Set("error", err.Error())
	q.Set("error_description", description)

	uri.RawQuery = q.Encode()
	return uri
}
