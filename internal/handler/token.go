package handler

import (
	"encoding/json"
	"fmt"
	"github.com/yael-castro/godi/internal/business"
	"github.com/yael-castro/godi/internal/model"
	"log"
	"mime"
	"net/http"
	"net/url"
)

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
			log.Println("ERROR", err)
			if err == nil {
				redirectURL = redirect
			}
		}

		exchange := model.Exchange{
			GrantType:         r.Form.Get("grant_type"),
			ClientId:          r.Form.Get("client_id"),
			AuthorizationCode: model.AuthorizationCode(r.Form.Get("code")),
			CodeVerifier:      model.CodeVerifier(r.Form.Get("code_verifier")),
			State:             model.State(r.Form.Get("state")),
			RedirectURL:       redirectURL,
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
