package business

import (
	"errors"
	"github.com/yael-castro/goauth/internal/model"
	"github.com/yael-castro/goauth/internal/repository"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

// authenticationTestCase common test case for the Authenticator interface
type authenticationTestCase struct {
	input       interface{}
	expectedErr error
}

// TestAuthenticator_Authenticate
// Test the method Authenticate for different implementations of the Authenticator interface
// The implementations that are tested
func TestAuthenticator_Authenticate(t *testing.T) {
	tdt := []struct {
		authenticator Authenticator
		tests         []authenticationTestCase
	}{
		{
			authenticator: OwnerAuthenticator{
				Obtainer: &repository.MockStorage[string, model.Owner]{
					"contacto@yael-castro.com": model.Owner{
						Id:       "contacto@yael-castro.com",
						Password: "$2a$10$g141w.TTnp5Bm/rLNqRRRevOSFhKBdV5KaJYxEDi9U5R9TgkZbfne", // yael.castro
					},
				},
			},
			tests: []authenticationTestCase{
				{
					input: model.Owner{
						Id:       "contacto@yael-castro.com",
						Password: "yael.castro",
					},
				},
				{
					input:       model.Owner{},
					expectedErr: model.AccessDenied,
				},
			},
		},
		{
			authenticator: ClientAuthenticator{
				Obtainer: repository.ObtainerFunc[string, model.Client](repository.MockStorage[string, model.Client]{
					"mobile": model.Client{
						AllowedOrigins: []string{"https://goauth.com"},
					},
				}.Obtain),
			},
			tests: []authenticationTestCase{
				{
					input: model.Application{
						Id: "mobile",
						RedirectURL: func() *url.URL {
							uri, _ := url.Parse("https://goauth.com")
							return uri
						}(),
					},
				},
				{
					input:       model.Application{},
					expectedErr: model.UnauthorizedClient,
				},
			},
		},
	}

	for _, v := range tdt {
		authenticator := v.authenticator
		t.Run(reflect.TypeOf(authenticator).String(), func(t *testing.T) {
			for i, v := range v.tests {
				t.Run(strconv.Itoa(i+1), func(t *testing.T) {
					err := authenticator.Authenticate(v.input)
					if !errors.Is(err, v.expectedErr) {
						t.Fatalf(`expected error "%v" but got "%v"`, v.expectedErr, err)
					}

					t.Log(err)
				})
			}
		})
	}
}
