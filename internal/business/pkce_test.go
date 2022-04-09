package business

import (
	"fmt"
	"github.com/yael-castro/godi/internal/model"
	"github.com/yael-castro/godi/internal/repository"
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
		expectedURL *url.URL
		expectedErr string
	}{
		// Success test case
		{
			input: model.Authorization{
				State:               "AAA",
				ClientId:            "a06a0630-31f5-4cc3-8e47-ea61a60c1199",
				ResponseType:        "code",
				CodeChallenge:       "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_~BCDEE",
				CodeChallengeMethod: "PLAIN",
				Scope:               "http://localhost/private",
				RedirectURL: func() *url.URL {
					uri, _ := url.Parse("http://localhost/callback")
					return uri
				}(),
				BasicAuth: model.Owner{
					Id:       "contacto@yael-castro.com",
					Password: "yael.castro",
				},
			},
			expectedURL: func() *url.URL {
				uri, _ := url.Parse("http://localhost/callback?code=ABC&state=AAA")
				return uri
			}(),
		},
		// Test case for invalid state
		{
			input: model.Authorization{
				ClientId: "a06a0630-31f5-4cc3-8e47-ea61a60c1199",
				RedirectURL: func() *url.URL {
					uri, _ := url.Parse("http://localhost/callback")
					return uri
				}(),
				BasicAuth: model.Owner{
					Id:       "contacto@yael-castro.com",
					Password: "yael.castro",
				},
			},
			expectedErr: model.InvalidRequest.Error(),
		},
		// Test case for invalid client id
		{
			input: model.Authorization{
				State:       "BBB",
				RedirectURL: &url.URL{},
			},
			expectedErr: model.UnauthorizedClient.Error(),
		},
		// Test case for invalid redirect url
		{
			input: model.Authorization{
				State:    "CCC",
				ClientId: "a06a0630-31f5-4cc3-8e47-ea61a60c1199",
				BasicAuth: model.Owner{
					Id:       "contacto@yael-castro.com",
					Password: "yael.castro",
				},
				RedirectURL: &url.URL{},
			},
			expectedErr: model.UnauthorizedClient.Error(),
		},
		// Invalid code_challenge_method
		{
			input: model.Authorization{
				State:    "DDD",
				ClientId: "a06a0630-31f5-4cc3-8e47-ea61a60c1199",
				BasicAuth: model.Owner{
					Id:       "contacto@yael-castro.com",
					Password: "yael.castro",
				},
				RedirectURL: func() *url.URL {
					uri, _ := url.Parse("http://localhost/callback")
					return uri
				}(),
			},
			expectedErr: model.InvalidRequest.Error(),
		},
		// Invalid code_challenge
		{
			input: model.Authorization{
				ClientId: "a06a0630-31f5-4cc3-8e47-ea61a60c1199",
				BasicAuth: model.Owner{
					Id:       "contacto@yael-castro.com",
					Password: "yael.castro",
				},
				State: "EEE",
				RedirectURL: func() *url.URL {
					uri, _ := url.Parse("http://localhost/callback")
					return uri
				}(),
				CodeChallengeMethod: "S256",
			},
			expectedErr: model.InvalidRequest.Error(),
		},
	}

	authorizer := AuthorizationCodeGrant{
		PKCE: ProofKeyCodeExchange{},
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
		CodeGenerator: CodeGeneratorFunc(func() model.AuthorizationCode {
			return "ABC"
		}),
		Storage: &repository.MockStorage{},
		Authenticator: BCryptAuthenticator{
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
			uri := authorizer.Authorize(v.input)

			err := uri.Query().Get("error")
			errDescription := uri.Query().Get("error_description")

			if err != fmt.Sprint(v.expectedErr) {
				t.Error(errDescription)
				t.Fatalf(`expected error "%v" got "%v"`, v.expectedErr, err)
			}

			if err != "" {
				t.Skipf("%v => %v", err, errDescription)
			}

			if uri.String() != v.expectedURL.String() {
				t.Fatalf(`unexpected redirect uri "%v"`, uri)
			}

			t.Log(uri)
		})
	}
}
