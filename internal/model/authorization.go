package model

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
		State               string
		CodeChallenge       string
		CodeChallengeMethod string
		// RedirectUri is not required by the spec, but your service should require it.
		// This URL must match one of the URLs the developer registered when creating the application,
		// and the authorization server should reject the request if it does not match (Optional)
		RedirectUri string
	}

	Application struct {
		ClientId     string
		ClientSecret string
		// AllowedOrigins valid redirect URIs
		AllowedOrigins []string
	}
)
