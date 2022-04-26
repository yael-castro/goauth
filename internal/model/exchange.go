package model

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/golang-jwt/jwt"
	"net/url"
	"strings"
	"time"
)

// Exchange contains the required data to exchange an authorization code for a token
type Exchange struct {
	GrantType string
	Application
	AuthorizationCode
	CodeVerifier
	State
	RedirectURL *url.URL
	// Session metadata of client
	// Is NOT part of the OAuth 2.0 protocol
	Session
}

// Token standard structure used when a token is returned
// (Contains the token and token metadata)
type Token struct {
	// Type indicates the token type
	//
	// Example: Bearer or Basic
	Type string `json:"type,omitempty"`
	// AccessToken: is the token returned
	AccessToken string `json:"access_token,omitempty"`
	// Scope represents the permissions that the token have
	Scope interface{} `json:"scope,omitempty"`
	// ExpiresIn token lifetime
	ExpiresIn *time.Duration `json:"expires_in,omitempty"`
}

// CodeVerifier is the code which the CodeChallenge is generated
type CodeVerifier string

// IsValid receives a code and method to check if the CodeVerifier match to the CodeChallengeMethod
func (v CodeVerifier) IsValid(code CodeChallenge, method CodeChallengeMethod) bool {
	method = CodeChallengeMethod(strings.ToUpper(string(method)))

	if method == "PLAIN" {
		return code == CodeChallenge(v)
	}

	hash := sha256.New().Sum([]byte(code))

	challenge := base64.URLEncoding.EncodeToString(hash)

	return challenge == string(code)
}

// Session details about the owner session
type Session struct {
	// IP address v4
	IP
	// Owner who is owner of this session
	Owner
	// TokenId is the token identifier like JTI
	TokenId string
	// UserAgent device
	UserAgent string
	// Expiration Token Lifetime
	Expiration time.Duration
}

// StandardClaims alias for jwt.StandardClaims
type StandardClaims = jwt.StandardClaims

// JWT JSON Web Token
type JWT struct {
	StandardClaims
	// Scope indicates the permissions that the JWT has
	Scope interface{} `json:"scp"`
}
