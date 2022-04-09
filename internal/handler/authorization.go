package handler

import (
	"fmt"
	"mime"
	"net/http"
	"net/url"

	"github.com/yael-castro/godi/internal/business"
	"github.com/yael-castro/godi/internal/model"
)

// NewAuthorizationHandler creates a http.HandleFunc using a business.Authorizer to handle authorization requests in
// the Authorization Code Grant flow described in the OAuth 2.0 protocol
func NewAuthorizationHandler(authorizer business.Authorizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
			return
		}

		if media != "application/x-www-form-urlencoded" {
			http.Error(w, fmt.Sprintf(`media "%s" is not supported`, media), http.StatusUnsupportedMediaType)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
			return
		}

		redirectURL := &[]url.URL{*r.URL}[0]

		if r.Form.Get("redirect_url") != "" {
			redirect, err := url.Parse(r.Form.Get("redirect_url"))
			if err == nil {
				redirectURL = redirect
			}
		}

		username, password, _ := r.BasicAuth()

		redirectURL = authorizer.Authorize(model.Authorization{
			ClientId:            r.Form.Get("client_id"),
			ClientSecret:        r.Form.Get("client_secret"),
			State:               model.State(r.Form.Get("state")),
			CodeChallenge:       model.CodeChallenge(r.Form.Get("code_challenge")),
			CodeChallengeMethod: model.CodeChallengeMethod(r.Form.Get("code_challenge_method")),
			Scope:               r.Form.Get("scope"),
			ResponseType:        r.Form.Get("response_type"),
			RedirectURL:         redirectURL,
			BasicAuth: model.Owner{
				Id:       username,
				Password: password,
			},
		})

		http.Redirect(w, r, redirectURL.String(), http.StatusFound)
	}
}
