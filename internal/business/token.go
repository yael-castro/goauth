package business

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
	"github.com/yael-castro/goauth/internal/model"
)

// TokenGenerator defines a provider of token generated from some data
type TokenGenerator interface {
	// GenerateToken generates a model.Token based on the received data
	GenerateToken(interface{}) (model.Token, error)
}

// _ "implement" constraint for JWTGenerator
var _ TokenGenerator = (*JWTGenerator)(nil)

type JWTGenerator struct {
	privateKey *rsa.PrivateKey
}

// SetPrivateKey parse the slice of bytes to a *rsa.PrivateKey
func (g *JWTGenerator) SetPrivateKey(privateKey []byte) (err error) {
	g.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	return
}

// GenerateToken generates a JWT based on the model.JWT received as parameter
func (g JWTGenerator) GenerateToken(i interface{}) (model.Token, error) {
	claims := i.(model.JWT)

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(g.privateKey)

	tkn := model.Token{
		Type:        "Bearer",
		AccessToken: token,
		Scope:       claims.Scope,
	}

	return tkn, err
}
