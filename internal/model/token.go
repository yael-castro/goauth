package model

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"time"
)

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
	// UserAgent device
	UserAgent string
	// Expiration Token Lifetime
	Expiration time.Duration
}
