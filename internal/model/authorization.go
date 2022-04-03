package model

import (
	"net/url"
	"path"
	"regexp"
	"strings"
)

// Authorization request of authorization following the protocol OAuth 2.0
type Authorization struct {
	// ClientId public application id
	ClientId string `json:"clientId,omitempty"`
	// ClientSecret private application secret
	ClientSecret string `json:"clientSecret,omitempty"`
	// Scope one or more scope values indicating additional access requested by the application (Optional)
	Scope string `json:"scope,omitempty"`
	// ResponseType expected response type (code, ...)
	ResponseType        string `json:"responseType,omitempty"`
	State               `json:"state,omitempty"`
	CodeChallenge       `json:"codeChallenge,omitempty"`
	CodeChallengeMethod `json:"codeChallengeMethod,omitempty"`
	// RedirectURL is not required by the spec, but your service should require it.
	// This URL must match one of the URLs the developer registered when creating the application,
	// and the authorization server should reject the request if it does not match (Optional)
	RedirectURL *url.URL `json:"redirectURL,omitempty"`
}

// Client defines an allowed client to make request for the Authorization Server
type Client struct {
	// Id public client identifier
	Id string
	// Secret optional secret client
	Secret string
	// AllowedOrigins origins to which the client can be redirected
	AllowedOrigins []string
}

// IsValidOrigin checks if the origin received as parameter match with some valid origin in AllowedOrigins
func (c Client) IsValidOrigin(origin string) bool {
	origin = path.Join(origin)

	for _, allowedOrigin := range c.AllowedOrigins {
		allowedOrigin = path.Join(allowedOrigin)

		if allowedOrigin == origin {
			return true
		}
	}

	return false
}

// State is used by the application to store request-specific data and/or prevent CSRF attacks (Recommended)
type State string

// IsValid indicates if the State is valid
func (s State) IsValid() (ok bool) {
	ok, _ = regexp.MatchString(`^[\x20-\x7E]+$`, string(s))
	return
}

// CodeChallenge code based on the CodeVerifier
type CodeChallenge string

// IsValid check if the CodeChallenge is valid
func (c CodeChallenge) IsValid() (ok bool) {
	ok, _ = regexp.MatchString(`^([-A-Z.a-z0-9]|_|~){43,128}$`, string(c))
	return
}

// CodeChallengeMethod is the method that the token endpoint (authorization endpoint) MUST use to verify
// the "code_verifier"
type CodeChallengeMethod string

// IsValid indicates if the CodeChallengeMethod is valid
func (m CodeChallengeMethod) IsValid() bool {
	method := strings.ToLower(string(m))
	return method == "plain" || method == "s256"
}

type AuthorizationCode string
