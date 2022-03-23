package business

import (
	"net/url"

	"github.com/yael-castro/godi/internal/model"
)

// ProofKeyCodeExchange is the "Authorization Code Grant" flow with the extension "Proof Key for Code Exchange"
// for the OAuth 2.0 protocol
type ProofKeyCodeExchange struct{}

func (p ProofKeyCodeExchange) Authorize(auth model.Authorization) *url.URL {
	panic("implement me")
}
