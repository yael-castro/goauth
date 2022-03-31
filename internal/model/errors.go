package model

import "fmt"

// NotFound error caused by missing resource
type NotFound string

// Error returns the string value of NotFound
func (n NotFound) Error() string {
	return string(n)
}

// DuplicateRecord error caused by duplicate primary key or duplicate record
type DuplicateRecord string

// Error returns the string value of DuplicateRecord
func (d DuplicateRecord) Error() string {
	return string(d)
}

// _ "implement" constraint for OAuthError
var _ error = OAuthError(0)

type OAuthError uint

// Error parsing the integer value to
func (e OAuthError) Error() string {
	switch e {
	case InvalidRequest:
		return "invalid_request"
	case AccessDenied:
		return "access_denied"
	case UnauthorizedClient:
		return "unauthorized_client"
	case UnsupportedResponseType:
		return "unsupported_response_type"
	case InvalidScope:
		return "invalid_scope"
	case ServerError:
		return "server_error"
	case TemporarilyUnavailable:
		return "temporarily_unavailable"
	}

	panic(fmt.Sprintf(`value "%d" is not supported`, e))
}

// Supported values for OAuthError
const (
	// InvalidRequest the client is not allowed to request an authorization code using this method,
	// for example if a confidential client attempts to use the implicit grant type
	InvalidRequest OAuthError = iota
	// AccessDenied the user or authorization server denied the request
	AccessDenied
	// UnauthorizedClient the client is not allowed to request an authorization code using this method,
	// for example if a confidential client attempts to use the implicit grant type
	UnauthorizedClient
	// UnsupportedResponseType the server does not support obtaining an authorization code using this method,
	// for example if the authorization server never implemented the implicit grant type
	UnsupportedResponseType
	// InvalidScope the requested scope is invalid or unknown
	InvalidScope
	// ServerError instead of displaying a 500 Internal Server Error page to the user,
	// the server can redirect with this error code.
	ServerError
	// TemporarilyUnavailable if the server is undergoing maintenance, or is otherwise unavailable,
	// this error code can be returned instead of responding with a 503 Service Unavailable status code
	TemporarilyUnavailable
)
