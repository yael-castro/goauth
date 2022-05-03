package business

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/yael-castro/goauth/internal/model"
	"github.com/yael-castro/goauth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// Authenticator defines the process of confirming that something is who it says it is
type Authenticator[T any] interface {
	// Authenticate check if something is who it says it is
	Authenticate(T) error
}

// _ "implement" constraint for OwnerAuthenticator
var _ Authenticator[model.Owner] = OwnerAuthenticator{}

type OwnerAuthenticator struct {
	repository.Obtainer[string, model.Owner]
}

// Authenticate validates a model.Owner to check if the password match to hashed password in database obtained by the owner id
func (o OwnerAuthenticator) Authenticate(owner model.Owner) (err error) {
	savedOwner, err := o.Obtainer.Obtain(owner.Id)
	if _, ok := err.(model.NotFound); err == redis.Nil || ok {
		err = fmt.Errorf(`%w: owner "%s" does not exists`, model.AccessDenied, owner.Id)
	}

	if err != nil {
		return
	}

	hashedPassword := savedOwner.Password

	ok := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(owner.Password)) == nil
	if !ok {
		err = fmt.Errorf("%w: invalid owner", model.AccessDenied)
	}

	return
}

// _ "implement" constraint for ClientAuthenticator
var _ Authenticator[model.Application] = ClientAuthenticator{}

// ClientAuthenticator authenticates a model.Application
type ClientAuthenticator struct {
	repository.Obtainer[string, model.Client]
}

// Authenticate validates a model.Application to check if the client credentials and redirect url match
// to some record in database
func (c ClientAuthenticator) Authenticate(application model.Application) (err error) {
	savedClient, err := c.Obtainer.Obtain(application.Id)
	if _, ok := err.(model.NotFound); ok || err == redis.Nil {
		return fmt.Errorf(`%w: client "%s" does not exist`, model.UnauthorizedClient, application.Id)
	}

	if err != nil {
		return
	}

	if !savedClient.IsValidOrigin(application.RedirectURL.String()) {
		err = fmt.Errorf("%w: invalid redirect_uri", model.UnauthorizedClient)
	}

	return
}
