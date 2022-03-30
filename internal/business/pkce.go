package business

import (
	"github.com/yael-castro/godi/internal/model"
	"github.com/yael-castro/godi/internal/repository"
	"net/url"
	"regexp"
)

// _ "implement" constraint for ProofKeyCodeExchange
var _ Authorizer = (*AuthorizationCodeGrant)(nil)

// AuthorizationCodeGrant authorize the allowed redirect URLs
type AuthorizationCodeGrant struct {
	finder repository.ApplicationFinder
}

// SetFinder initializes the instance of repository.ApplicationFinder used to get the model.Application
func (c *AuthorizationCodeGrant) SetFinder(finder repository.ApplicationFinder) {
	c.finder = finder
}

// Authorize validate the application (model.Application) obtained with the received data (model.Authorization)
//
// In resume...
//
// 1. Identifies the app using the client id and client secret
//
// 2. Validates the received state
//
// 3. Validates the redirect uri
func (c AuthorizationCodeGrant) Authorize(auth model.Authorization) *url.URL {
	q := auth.RedirectURL.Query()

	app, err := c.finder.FindApplication(auth.ClientId)
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
	if ok, _ := regexp.MatchString(`[\x20-\x7E]`, auth.State); !ok {
		q.Set("error", model.InvalidRequest.Error())
		q.Set("error_description", "invalid state")
		goto end
	}

	// Identifying application using the client id and client secret
	if app.ClientId != auth.ClientId || app.ClientSecret != auth.ClientSecret {
		q.Set("error", model.UnauthorizedClient.Error())
		q.Set("error_description", "client credentials does not match")
		goto end
	}

	// Validation of redirect uri
	for _, origin := range app.AllowedOrigins {
		if auth.RedirectURL.RawPath == origin {
			return auth.RedirectURL
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
	repository.SessionStorage
}

// Authorize validates the code_challenge and code_challenge_method
func (p ProofKeyCodeExchange) Authorize(auth model.Authorization) *url.URL {
	if uri := p.Authorizer.Authorize(auth); uri.Query().Get("error") != "" { // Redirect URI validation
		return uri
	}

	q := auth.RedirectURL.Query()

	if auth.CodeChallengeMethod == "" {
		auth.CodeChallengeMethod = "PLAIN"
	}

	if auth.CodeChallengeMethod != "PLAIN" && auth.CodeChallengeMethod != "S256" {
		q.Set("error", model.InvalidRequest.Error())
		q.Set("error_description", "invalid code_challenge_method, must be PLAIN or S256")
		goto end
	}

	// TODO validate code challenge

	if err := p.SessionStorage.CreateSession(auth); err != nil {
		q.Set("error", model.ServerError.Error())
		q.Set("error_description", err.Error())
	}

end:
	auth.RedirectURL.RawQuery = q.Encode()
	return auth.RedirectURL
}
