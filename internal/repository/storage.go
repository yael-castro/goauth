package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/yael-castro/godi/internal/model"
	"time"
)

// Constants for Authorization Code Life Time
const (
	MaximumAuthorizationCodeLifeTime = 10 * time.Minute
	MediumAuthorizationCodeLifeTime  = MaximumAuthorizationCodeLifeTime / 2
)

// Storage defines a general store
// For example: potato storage or client storage
type Storage interface {
	// Create creates a record using the received data
	Create(string, interface{}) error
	// Obtain obtains record by id
	Obtain(string) (interface{}, error)
	// Delete removed record from the storage
	Delete(string) error
}

// _ "implement" constraint for StateStorage
var _ Storage = StateStorage{}

// StateStorage storage of states related to authorization requests
// Basically saves instances of model.Authorization
type StateStorage struct {
	*redis.Client
}

// authorizationKey creates a new key to be an identifier in the redis database based on the received state
func (StateStorage) authorizationKey(code model.AuthorizationCode) string {
	return "authorization:" + string(code)
}

// Create receives a model.Authorization and uses your data to create a state-identified record
//
// Notes:
//
// - The record only is created if the state does not exist
//
// - If the record exists an error of type model.DuplicateRecord is returned
func (s StateStorage) Create(code string, i interface{}) error {
	auth := i.(model.Authorization)

	cmd := s.SetNX(context.TODO(), s.authorizationKey(model.AuthorizationCode(code)), model.BinaryJSON{I: auth}, MediumAuthorizationCodeLifeTime)

	flag, err := cmd.Result()
	if err != nil {
		return err
	}

	if !flag {
		err = model.DuplicateRecord(fmt.Sprintf(`athorization code "%s" already exists`, code))
	}

	return err
}

// Obtain search a saved instance of model.Authorization by the state
func (s StateStorage) Obtain(code string) (interface{}, error) {
	cmd := s.Get(context.TODO(), s.authorizationKey(model.AuthorizationCode(code)))

	serializedData, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	auth := model.Authorization{}

	err = json.Unmarshal([]byte(serializedData), &auth)
	return auth, err
}

// Delete removes a record using the state received as parameter
//
// Note: if the record does not exist, it returns NO errors
func (s StateStorage) Delete(code string) error {
	return s.Del(context.TODO(), s.authorizationKey(model.AuthorizationCode(code))).Err()
}

// _ "implement" constraint for *MockStorage
var _ Storage = (*MockStorage)(nil)

// MockStorage store for model.Authorization
type MockStorage map[string]interface{}

// Create saves a model.Authorization in m
func (m *MockStorage) Create(code string, i interface{}) error {
	(*m)[code] = i
	return nil
}

// Obtain search a model.Authorization by state
func (m MockStorage) Obtain(code string) (interface{}, error) {
	return m[code], nil
}

// Delete removes a record by state
func (m *MockStorage) Delete(code string) error {
	delete(*m, code)
	return nil
}

// _ "implement" constraint for OwnerStorage
var _ Storage = OwnerStorage{}

// OwnerStorage is store for all data directly related to o
type OwnerStorage struct {
	*redis.Client
}

// ownerKey generates key to save the owner data
func (o OwnerStorage) ownerKey(ownerId string) string {
	return "owner:" + ownerId
}

// Create not implemented yet
func (o OwnerStorage) Create(ownerId string, owner interface{}) error {
	cmd := o.SetNX(context.TODO(), o.ownerKey(ownerId), owner.(model.Owner).Password, 0)

	return cmd.Err()
}

// Obtain search a model.Owner by ownerId (user id)
func (o OwnerStorage) Obtain(ownerId string) (i interface{}, err error) {
	owner := model.Owner{Id: ownerId}

	owner.Password, err = o.Get(context.TODO(), o.ownerKey(ownerId)).Result()
	if err != nil {
		return
	}

	i = owner
	return
}

// Delete deletes a model.Owner by ownerId
func (o OwnerStorage) Delete(ownerId string) error {
	cmd := o.Del(context.TODO(), o.ownerKey(ownerId))
	return cmd.Err()
}

// _ "implement" constraint for SessionStorage
var _ Storage = SessionStorage{}

// SessionStorage storage for owner enabled sessions
type SessionStorage struct {
	*redis.Client
}

// sessionKey creates a session key based on the token id
func (s SessionStorage) sessionKey(tokenId string) string {
	return "session:" + tokenId
}

// Create creates a record of model.Session with the received tokenId to
func (s SessionStorage) Create(tokenId string, i interface{}) error {
	session := i.(model.Session)

	cmd := s.SetNX(context.TODO(), s.sessionKey(tokenId), model.BinaryJSON{I: session}, session.Expiration)

	wasCreated, err := cmd.Result()
	if err != nil {
		return err
	}

	if !wasCreated {
		err = model.DuplicateRecord(`session already exists`)
	}

	return err
}

// Obtain search enabled session by token id
func (s SessionStorage) Obtain(tokenId string) (i interface{}, err error) {
	serialized, err := s.Get(context.TODO(), s.sessionKey(tokenId)).Result()

	session := model.Session{}

	err = json.Unmarshal([]byte(serialized), &session)
	if err != nil {
		return
	}

	i = session

	return
}

// Delete revokes a session by token id
func (s SessionStorage) Delete(tokenId string) error {
	return s.Del(context.TODO(), s.sessionKey(tokenId)).Err()
}
