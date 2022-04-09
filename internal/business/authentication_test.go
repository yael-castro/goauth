package business

import (
	"github.com/yael-castro/godi/internal/model"
	"github.com/yael-castro/godi/internal/repository"
	"strconv"
	"testing"
)

// TestProofKeyCodeExchange_Authorize
// Test the implementation OwnerAuthenticator of Authenticator
func TestBCryptAuthenticator_Authenticate(t *testing.T) {
	tdt := []struct {
		input       interface{}
		output      bool
		expectedErr error
	}{
		{
			input: model.Owner{
				Id:       "contacto@yael-castro.com",
				Password: "yael.castro",
			},
			output: true,
		},
	}

	authenticator := OwnerAuthenticator{
		Storage: &repository.MockStorage{
			"contacto@yael-castro.com": model.Owner{
				Id:       "contacto@yael-castro.com",
				Password: "$2a$10$g141w.TTnp5Bm/rLNqRRRevOSFhKBdV5KaJYxEDi9U5R9TgkZbfne",
			},
		},
	}

	for i, v := range tdt {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ok, err := authenticator.Authenticate(v.input)
			if err != v.expectedErr {
				t.Fatalf(`expected error "%v" but got "%v"`, v.expectedErr, err)
			}

			if ok != v.output {
				t.Fatalf(`expected "%v" got "%v""`, v.output, ok)
			}

			t.Log(ok)
		})
	}
}
