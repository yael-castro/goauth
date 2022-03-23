package model

// Supported scopes (permissions)
const (
	// Undefined value that indicates "no permissions"
	Undefined Scope = 0
	Create    Scope = 1 << iota
	Read
	Update
	Delete
)

// Scope bit mask that indicates the different scopes (permissions) for the application
// Deprecated: at the moment is only a draft
type Scope uint

// Is validates if the scope received as parameter match to s
func (s Scope) Is(scope Scope) bool {
	return s&scope == scope
}

type (
	// Authorization request of authorization following the protocol OAuth 2.0
	Authorization struct {
		// ClientId public application id
		ClientId string
		// ClientSecret private application secret
		ClientSecret string
		// Scope one or more scope values indicating additional access requested by the application (Optional)
		Scope string
		// ResponseType expected response type (code, ...)
		ResponseType string
		// State is used by the application to store request-specific data and/or prevent CSRF attacks (Recommended)
		State         string
		CodeChallenge string
		// CodeChallengeMethod is the method that the token endpoint (authorization endpoint) MUST use to verify
		// the "code_verifier"
		CodeChallengeMethod string
		// RedirectUri is not required by the spec, but your service should require it.
		// This URL must match one of the URLs the developer registered when creating the application,
		// and the authorization server should reject the request if it does not match (Optional)
		RedirectUri string
	}

	Application struct {
		ClientId     stringkkk
		ClientSecret string
		// AllowedOrigins valid redirect URIs
		AllowedOrigins []string
	}
)
