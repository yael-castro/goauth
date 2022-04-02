package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/yael-castro/godi/internal/model"
	"time"
)

const MaximumAuthorizationCodeLifetime = 10 * time.Minute

// Storage defines a general store
// For example: potato storage or client session storage
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

	cmd := s.SetNX(context.TODO(), s.authorizationKey(model.AuthorizationCode(code)), model.BinaryJSON{I: auth}, MaximumAuthorizationCodeLifetime)

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

// _ "implement" constraint for *MockStateStore
var _ Storage = (*MockStateStore)(nil)

// MockStateStore store for model.Authorization
type MockStateStore map[model.AuthorizationCode]model.Authorization

// Create saves a model.Authorization in m
func (m *MockStateStore) Create(code string, i interface{}) error {
	auth := i.(model.Authorization)

	(*m)[model.AuthorizationCode(code)] = auth

	return nil
}

// Obtain search a model.Authorization by state
func (m MockStateStore) Obtain(code string) (interface{}, error) {
	return m[model.AuthorizationCode(code)], nil
}

// Delete removes a record by state
func (m *MockStateStore) Delete(code string) error {
	delete(*m, model.AuthorizationCode(code))
	return nil
}
