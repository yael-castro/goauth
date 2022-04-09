package business

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/yael-castro/godi/internal/model"
	"github.com/yael-castro/godi/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// Authenticator defines the process of confirming that something is who it says it is
type Authenticator interface {
	// Authenticate check if something is who it says it is
	Authenticate(interface{}) (bool, error)
}

// _ "implement" constraint for BCryptAuthenticator
var _ Authenticator = BCryptAuthenticator{}

type BCryptAuthenticator struct {
	repository.Storage
}

// Authenticate validates a model.Owner to check if the password match to hashed password in database obtained by the owner id
func (p BCryptAuthenticator) Authenticate(i interface{}) (ok bool, err error) {
	owner := i.(model.Owner)

	ownerData, err := p.Storage.Obtain(owner.Id)
	if _, ok := err.(model.NotFound); errors.Is(err, redis.Nil) || ok {
		return false, nil
	}

	if err != nil {
		return
	}

	hashedPassword := ownerData.(model.Owner).Password

	ok = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(owner.Password)) == nil
	return
}
