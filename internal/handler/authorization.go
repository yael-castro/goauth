package handler

import (
	"net/http"
	"net/url"

	"github.com/yael-castro/godi/internal/business"
	"github.com/yael-castro/godi/internal/model"
)

// NewAuthorizationHandler creates a http.HandleFunc using a business.Authorizer to handle authorization requests in
// the Authorization Code Grant flow described in the OAuth 2.0 protocol
func NewAuthorizationHandler(authorizer business.Authorizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectURL := &[]url.URL{*r.URL}[0]

		if r.Method != http.MethodGet {
			business.OAuthError(redirectURL, model.InvalidRequest, "method not allowed")
			goto end
		}

		if err := r.ParseForm(); err != nil {
			business.OAuthError(redirectURL, model.InvalidRequest, err.Error())
			goto end
		}

		if r.Form.Get("redirect_url") != "" {
			redirect, err := url.Parse(r.Form.Get("redirect_url"))
			if err == nil {
				redirectURL = redirect
			}
		}

		redirectURL = authorizer.Authorize(model.Authorization{
			ClientId:            r.Form.Get("client_id"),
			ClientSecret:        r.Form.Get("client_secret"),
			State:               model.State(r.Form.Get("state")),
			CodeChallenge:       model.CodeChallenge(r.Form.Get("code_challenge")),
			CodeChallengeMethod: model.CodeChallengeMethod(r.Form.Get("code_challenge_method")),
			Scope:               r.Form.Get("scope"),
			ResponseType:        r.Form.Get("response_type"),
			RedirectURL:         redirectURL,
		})

	end:
		http.Redirect(w, r, redirectURL.String(), http.StatusFound)
	}
}
