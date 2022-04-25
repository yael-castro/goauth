package business

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/yael-castro/goauth/internal/model"
	"github.com/yael-castro/goauth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// Authenticator defines the process of confirming that something is who it says it is
type Authenticator interface {
	// Authenticate check if something is who it says it is
	Authenticate(interface{}) error
}

// _ "implement" constraint for OwnerAuthenticator
var _ Authenticator = OwnerAuthenticator{}

type OwnerAuthenticator struct {
	repository.Storage
}

// Authenticate validates a model.Owner to check if the password match to hashed password in database obtained by the owner id
func (o OwnerAuthenticator) Authenticate(i interface{}) (err error) {
	owner := i.(model.Owner)

	ownerData, err := o.Storage.Obtain(owner.Id)
	if _, ok := err.(model.NotFound); err == redis.Nil || ok {
		err = fmt.Errorf(`%w: owner "%s" does not exists`, model.AccessDenied, owner.Id)
	}

	if err != nil {
		return
	}

	hashedPassword := ownerData.(model.Owner).Password

	ok := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(owner.Password)) == nil
	if !ok {
		err = fmt.Errorf("%w: invalid owner", model.AccessDenied)
	}

	return
}

// _ "implement" constraint for ClientAuthenticator
var _ Authenticator = ClientAuthenticator{}

// ClientAuthenticator authenticates a model.Application
type ClientAuthenticator struct {
	repository.Finder
}

// Authenticate validates a model.Application to check if the client credentials and redirect url match
// to some record in database
func (c ClientAuthenticator) Authenticate(i interface{}) (err error) {
	application := i.(model.Application)

	data, err := c.Finder.Find(application.Id)
	if _, ok := err.(model.NotFound); ok || err == redis.Nil {
		return fmt.Errorf(`%w: client "%s" does not exist`, model.UnauthorizedClient, application.Id)
	}

	if err != nil {
		return
	}

	savedClient := data.(model.Client)

	if !savedClient.IsValidOrigin(application.RedirectURL.String()) {
		err = fmt.Errorf("%w: invalid redirect_uri", model.UnauthorizedClient)
	}

	return
}
