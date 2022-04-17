package business

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/yael-castro/godi/internal/model"
	"github.com/yael-castro/godi/internal/repository"
	"path"
	"time"
)

// Authorizer handles authorization requests for any flow of the OAuth 2.0 protocol
type Authorizer interface {
	// Authorize receives an authorization request and returns the redirect uri
	Authorize(model.Authorization) (model.AuthorizationCode, error)
}

// CodeExchanger defines the change of authorization code for a token
type CodeExchanger interface {
	// ExchangeCode receives a model.Exchange to change for some token
	ExchangeCode(model.Exchange) (model.Token, error)
}

// CodeGrant defines the interface related to the authorization code grant
type CodeGrant interface {
	Authorizer
	CodeExchanger
}

// _ "implement" constraints for ProofKeyCodeExchange
var _ CodeGrant = (*AuthorizationCodeGrant)(nil)

// AuthorizationCodeGrant made the validations that correspond to the Authorization Code Grant flow
type AuthorizationCodeGrant struct {
	// PKCE must be an implementation of the Proof Key for Code Exchange extension
	PKCE CodeChallengeValidator
	// ScopeParser parses a scope from string
	ScopeParser
	// CodeGenerator generate a code to save and validate the requests
	CodeGenerator
	TokenGenerator
	Owner  Authenticator
	Client Authenticator
	// CodeStorage store for all exchange codes generated
	CodeStorage repository.Storage
	// SessionStorage store for all exchange codes
	SessionStorage repository.Storage
}

// ExchangeCode using the model.Exchange search a record of mode.Authorization using the model.AuthorizationCode
// then if the record exists, the model.State, model.CodeChallenge
func (c AuthorizationCodeGrant) ExchangeCode(exchange model.Exchange) (tkn model.Token, err error) {
	// TODO validates the grant_type
	i, err := c.CodeStorage.Obtain(string(exchange.AuthorizationCode))
	if err != nil {
		return
	}

	authorization := i.(model.Authorization)

	if authorization.State != exchange.State {
		return model.Token{}, fmt.Errorf("%w: state does not match", model.AccessDenied)
	}

	if c.PKCE != nil {
		if !exchange.CodeVerifier.IsValid(authorization.CodeChallenge, authorization.CodeChallengeMethod) {
			return model.Token{}, fmt.Errorf("%w: code_verifier does not match to code_challenge", model.AccessDenied)
		}
	}

	if authorization.Application.Id != exchange.ClientId {
		return model.Token{}, errors.New("client_id does not match")
	}

	if path.Join(authorization.RedirectURL.String()) != path.Join(exchange.RedirectURL.String()) {
		return model.Token{}, fmt.Errorf("%w: redirect_uri does not match to the first redirect_uri", model.AccessDenied)
	}

	scope, err := c.ParseScope(authorization.Scope)
	if err != nil {
		return
	}

	token := model.JWT{
		// UserId: authorization.BasicAuth.Id,
		Scope: scope,
		StandardClaims: model.StandardClaims{
			Id:       uuid.New().String(),
			Issuer:   "go-auth",
			Subject:  authorization.BasicAuth.Id,
			Audience: authorization.Application.Id,
			IssuedAt: time.Now().Unix(),
		},
	}

	tkn, err = c.GenerateToken(token)
	if err != nil {
		return
	}

	// TODO check the data saved using the session storage
	err = c.SessionStorage.Create(token.Id, exchange.Session)
	if err != nil {
		return
	}

	err = c.CodeStorage.Delete(string(exchange.AuthorizationCode))
	return
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
//
// 4. Saves the session of this authorization request using the random code generated by
// the CodeGenerator
func (c AuthorizationCodeGrant) Authorize(a model.Authorization) (code model.AuthorizationCode, err error) {
	if a.ResponseType != "code" {
		return "", fmt.Errorf(`%w: "%s" is not supported`, model.UnsupportedResponseType, a.ResponseType)
	}

	err = c.Client.Authenticate(a.Application)
	if err != nil {
		return // model.FailedAuthentication
	}

	err = c.Owner.Authenticate(a.BasicAuth)
	if err != nil {
		return // model.FailedAuthentication
	}

	if !a.State.IsValid() {
		return "", fmt.Errorf("%w: state is not valid", model.InvalidRequest)
	}

	// Proof Key for Code Exchange (Extension)
	if c.PKCE != nil {
		if err = c.PKCE.ValidateCodeChallenge(a.CodeChallenge, a.CodeChallengeMethod); err != nil {
			return // Validation error
		}
	}

	_, err = c.ParseScope(a.Scope)
	if err != nil {
		return // Invalid scope
	}

	code = c.GenerateCode()
	if err = c.CodeStorage.Create(string(code), a); err != nil {
		return // Server error
	}

	return
}

// CodeChallengeValidator defines the additional validation required in the PKCE extension
type CodeChallengeValidator interface {
	// ValidateCodeChallenge validates the code_challenge and code_challenge_method as part of the PKCE extension
	ValidateCodeChallenge(model.CodeChallenge, model.CodeChallengeMethod) error
}

// ProofKeyCodeExchange is the "Authorization Code Grant" flow with the extension "Proof Key for Code Exchange"
// for the OAuth 2.0 protocol
type ProofKeyCodeExchange struct{}

// ValidateCodeChallenge validates the code_challenge and code_challenge_method
func (p ProofKeyCodeExchange) ValidateCodeChallenge(challenge model.CodeChallenge, method model.CodeChallengeMethod) error {
	if !method.IsValid() {
		return fmt.Errorf(`%w: invalid code_challenge_method "%s", must be PLAIN or S256`, model.InvalidRequest, challenge)
	}

	// TODO validate correctly the code challenge based on the model.CodeChallengeMethod
	// BUG: if code_challenge_method is S256 and code_challenge contains the ~ character
	// it is valid, which should be wrong because when using the S256 method the information
	// is sent in Base64 (URL) so it should not contain the ~ character at all
	//
	// This BUG should not cause any security problems but it does generate conflicts
	// when validating the code_verifier if the code_challenge contains the ~ character
	if !challenge.IsValid() {
		return fmt.Errorf("%w: %s", model.InvalidRequest, "invalid code_challenge")
	}

	return nil
}
