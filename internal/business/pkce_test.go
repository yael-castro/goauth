package business

import (
	"errors"
	"github.com/yael-castro/goauth/internal/model"
	"github.com/yael-castro/goauth/internal/repository"
	"net/url"
	"strconv"
	"testing"
)

// TestProofKeyCodeExchange_Authorize
//
// This tests the authorization request in the "Authorization Code Grant" flow with the "Proof Key for Code Exchange"
// extension of the OAuth 2.0 protocol using the ProofKeyCodeExchange and AuthorizationCodeGrant implementations
// of Authorizer with mock implementations for storage
//
func TestProofKeyCodeExchange_Authorize(t *testing.T) {
	tdt := []struct {
		input       model.Authorization
		expectedErr error
	}{
		// Success test case
		{
			input: model.Authorization{
				ResponseType: "code",
				Application: model.Application{
					Id: "a06a0630-31f5-4cc3-8e47-ea61a60c1199",
					RedirectURL: func() *url.URL {
						uri, _ := url.Parse("http://localhost/callback")
						return uri
					}(),
				},
				BasicAuth: model.Owner{
					Id:       "contacto@yael-castro.com",
					Password: "yael.castro",
				},
				State:               "AAA",
				Scope:               "read:ff write:afa",
				CodeChallenge:       "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_~BCDEE",
				CodeChallengeMethod: "PLAIN",
			},
		},
		// Test case for invalid state
		{
			input: model.Authorization{
				ResponseType: "code",
				Application: model.Application{
					Id: "a06a0630-31f5-4cc3-8e47-ea61a60c1199",
					RedirectURL: func() *url.URL {
						uri, _ := url.Parse("http://localhost/callback")
						return uri
					}(),
				},
				BasicAuth: model.Owner{
					Id:       "contacto@yael-castro.com",
					Password: "yael.castro",
				},
			},
			expectedErr: model.InvalidRequest,
		},
		// Test case for invalid client id
		{
			input: model.Authorization{
				ResponseType: "code",
				Application: model.Application{
					RedirectURL: &url.URL{},
				},
				State: "BBB",
			},
			expectedErr: model.UnauthorizedClient,
		},
		// Test case for invalid redirect url
		{
			input: model.Authorization{
				ResponseType: "code",
				Application: model.Application{
					Id:          "a06a0630-31f5-4cc3-8e47-ea61a60c1199",
					RedirectURL: &url.URL{},
				},
				BasicAuth: model.Owner{
					Id:       "contacto@yael-castro.com",
					Password: "yael.castro",
				},
				State: "CCC",
			},
			expectedErr: model.UnauthorizedClient,
		},
		// Invalid code_challenge_method
		{
			input: model.Authorization{
				ResponseType: "code",
				Application: model.Application{
					Id: "a06a0630-31f5-4cc3-8e47-ea61a60c1199",
					RedirectURL: func() *url.URL {
						uri, _ := url.Parse("http://localhost/callback")
						return uri
					}(),
				},
				BasicAuth: model.Owner{
					Id:       "contacto@yael-castro.com",
					Password: "yael.castro",
				},
				State: "DDD",
			},
			expectedErr: model.InvalidRequest,
		},
		// Invalid code_challenge
		{
			input: model.Authorization{
				ResponseType: "code",
				Application: model.Application{
					Id: "a06a0630-31f5-4cc3-8e47-ea61a60c1199",
					RedirectURL: func() *url.URL {
						uri, _ := url.Parse("http://localhost/callback")
						return uri
					}(),
				},
				BasicAuth: model.Owner{
					Id:       "contacto@yael-castro.com",
					Password: "yael.castro",
				},
				State:               "EEE",
				CodeChallengeMethod: "S256",
			},
			expectedErr: model.InvalidRequest,
		},
	}

	authorizer := AuthorizationCodeGrant{
		PKCE: ProofKeyCodeExchange{},
		Client: ClientAuthenticator{
			Finder: repository.MockClientFinder{
				"a06a0630-31f5-4cc3-8e47-ea61a60c1199": {
					Id:             "a06a0630-31f5-4cc3-8e47-ea61a60c1199",
					AllowedOrigins: []string{"http://localhost/callback/", "http://localhost:8080/callback/"},
				},
				"4cc3-8e47-ea61a60c1199-a06a0630-31f5": {
					Id:             "4cc3-8e47-ea61a60c1199-a06a0630-31f5",
					AllowedOrigins: []string{"http://localhost/callback/", "http://localhost:8080/callback/"},
				},
			},
		},
		CodeGenerator:  GenerateRandomCode,
		CodeStorage:    &repository.MockStorage{},
		SessionStorage: &repository.MockStorage{},
		ScopeParser:    NewScopeParser(),
		Owner: OwnerAuthenticator{
			Storage: &repository.MockStorage{
				"contacto@yael-castro.com": model.Owner{
					Id:       "contacto@yael-castro.com",
					Password: "$2a$10$g141w.TTnp5Bm/rLNqRRRevOSFhKBdV5KaJYxEDi9U5R9TgkZbfne",
				},
			},
		},
	}

	for i, v := range tdt {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			code, err := authorizer.Authorize(v.input)

			if !errors.Is(err, v.expectedErr) {
				t.Fatalf(`expected error "%v" got "%v"`, v.expectedErr, err)
			}

			if err != nil {
				t.Skipf("%v => %v", errors.Unwrap(err), err)
			}

			t.Log(code)
		})
	}
}
