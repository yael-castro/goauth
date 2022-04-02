package business

import (
	"github.com/yael-castro/godi/internal/model"
	"github.com/yael-castro/godi/internal/repository"
	"net/url"
)

// _ "implement" constraint for ProofKeyCodeExchange
var _ Authorizer = (*AuthorizationCodeGrant)(nil)

// AuthorizationCodeGrant defines the Authorization Code Grant
type AuthorizationCodeGrant struct {
	finder repository.Finder
}

// SetFinder initializes the instance of repository.ClientFinder used to get the model.Client
func (c *AuthorizationCodeGrant) SetFinder(finder repository.ClientFinder) {
	c.finder = finder
}

// Authorize validate the client model.Client obtained with the received data (model.Authorization)
//
// In resume...
//
// 1. Identifies the client using the client id and client secret
//
// 2. Validates the received state
//
// 3. Validates the redirect uri
func (c AuthorizationCodeGrant) Authorize(auth model.Authorization) *url.URL {
	q := auth.RedirectURL.Query()

	var client model.Client

	i, err := c.finder.Find(auth.ClientId)
	// Obtaining the application data by client id
	if _, ok := err.(model.NotFound); ok {
		q.Set("error", model.UnauthorizedClient.Error())
		q.Set("error_description", err.Error())
		goto end
	}

	// Handling internal server errors
	if err != nil {
		q.Set("error", model.ServerError.Error())
		q.Set("error_description", err.Error())
		goto end
	}

	// Validating state
	if auth.State.IsValid() {
		q.Set("error", model.InvalidRequest.Error())
		q.Set("error_description", "invalid state")
		goto end
	}

	client = i.(model.Client)

	// Identifying client using the client id and client secret
	if client.Id != auth.ClientId || client.Secret != auth.ClientSecret {
		q.Set("error", model.UnauthorizedClient.Error())
		q.Set("error_description", "client credentials does not match")
		goto end
	}

	// Validation of redirect uri
	for _, origin := range client.AllowedOrigins {
		if auth.RedirectURL.RawPath == origin {
			goto end
		}
	}

	q.Set("error", model.UnauthorizedClient.Error())
	q.Set("error_description", "invalid redirect uri")

end:
	auth.RedirectURL.RawQuery = q.Encode()
	return auth.RedirectURL
}

// _ "implement" constraint for ProofKeyCodeExchange
var _ Authorizer = ProofKeyCodeExchange{}

// ProofKeyCodeExchange is the "Authorization Code Grant" flow with the extension "Proof Key for Code Exchange"
// for the OAuth 2.0 protocol
type ProofKeyCodeExchange struct {
	// Authorizer must be an implementation of Authorize Code Grant Flow
	Authorizer
	repository.Storage
}

// Authorize validates the code_challenge and code_challenge_method and, if all is ok,
// saves the session of this authorization requests using the received state
func (p ProofKeyCodeExchange) Authorize(auth model.Authorization) *url.URL {
	if uri := p.Authorizer.Authorize(auth); uri.Query().Get("error") != "" { // Redirect URI validation
		return uri
	}

	q := auth.RedirectURL.Query()

	if auth.CodeChallengeMethod.IsValid() {
		q.Set("error", model.InvalidRequest.Error())
		q.Set("error_description", "invalid code_challenge_method, must be PLAIN or S256")
		goto end
	}

	if auth.CodeChallenge.IsValid() {
		q.Set("error", model.InvalidRequest.Error())
		q.Set("error_description", "invalid code_challenge")
		goto end
	}

	if err := p.Storage.Create(auth); err != nil {
		if _, ok := err.(model.DuplicateRecord); ok {
			q.Set("error", model.InvalidRequest.Error())
			q.Set("error_description", err.Error())
			goto end
		}

		q.Set("error", model.ServerError.Error())
		q.Set("error_description", err.Error())
	}

end:
	auth.RedirectURL.RawQuery = q.Encode()
	return auth.RedirectURL
}
