package handler

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/url"

	"github.com/yael-castro/goauth/internal/business"
	"github.com/yael-castro/goauth/internal/model"
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

		//if media == "" {
		// TODO do content negotiation to render the page?
		//}

		if media != "application/x-www-form-urlencoded" {
			http.Error(w, fmt.Sprintf(`media "%s" is not supported`, media), http.StatusUnsupportedMediaType)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		redirectURL := &[]url.URL{*r.URL}[0]

		if r.Form.Get("redirect_uri") != "" {
			redirect, err := url.Parse(r.Form.Get("redirect_uri"))
			if err == nil {
				redirectURL = redirect
			}
		}

		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", "Basic")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		a := model.Authorization{
			Application: model.Application{
				Id:          r.Form.Get("client_id"),
				Secret:      r.Form.Get("client_secret"),
				RedirectURL: redirectURL,
			},
			State:               model.State(r.Form.Get("state")),
			CodeChallenge:       model.CodeChallenge(r.Form.Get("code_challenge")),
			CodeChallengeMethod: model.CodeChallengeMethod(r.Form.Get("code_challenge_method")),
			Scope:               r.Form.Get("scope"),
			ResponseType:        r.Form.Get("response_type"),
			BasicAuth: model.Owner{
				Id:       username,
				Password: password,
			},
		}

		oauthErr := model.OAuthError(0)

		code, err := authorizer.Authorize(a)
		switch {
		case errors.As(err, &oauthErr):
			OAuthError(redirectURL, oauthErr, err.Error())

		case err != nil:
			OAuthError(redirectURL, model.ServerError, err.Error())

		default:
			redirectURL.RawQuery = url.Values{
				"code":  {string(code)},
				"state": {string(a.State)},
			}.Encode()
		}

		http.Redirect(w, r, redirectURL.String(), http.StatusFound)
	}
}
