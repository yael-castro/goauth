package business

import (
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
	uri, _ := url.Parse(auth.RedirectUri)

	flow, ok := o[auth.ResponseType]
	if !ok {
		q := uri.Query()

		q.Set("error", model.UnsupportedResponseType.Error())

		uri.RawQuery = q.Encode()
		return uri
	}

	return flow.Authorize(auth)
}
