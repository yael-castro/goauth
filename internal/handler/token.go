package handler

import (
	"encoding/json"
	"fmt"
	"github.com/yael-castro/goauth/internal/business"
	"github.com/yael-castro/goauth/internal/model"
	"mime"
	"net/http"
	"net/url"
)

// NewTokenHandler handler
func NewTokenHandler(exchanger business.CodeExchanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed)
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

		// TODO support port scanning
		ip, _ := model.NewIP(r.RemoteAddr)

		exchange := model.Exchange{
			GrantType: r.Form.Get("grant_type"),
			Application: model.Application{
				Id:          r.Form.Get("client_id"),
				Secret:      r.Form.Get("client_secret"),
				RedirectURL: redirectURL,
			},
			AuthorizationCode: model.AuthorizationCode(r.Form.Get("code")),
			CodeVerifier:      model.CodeVerifier(r.Form.Get("code_verifier")),
			State:             model.State(r.Form.Get("state")),
			Session: model.Session{
				UserAgent: r.UserAgent(),
				IP:        ip,
			},
		}

		token, err := exchanger.ExchangeCode(exchange)
		switch {
		case err != nil:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(token)
		}
	}
}
