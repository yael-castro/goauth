package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/yael-castro/godi/internal/model"
	"time"
)

// Storage defines a general storage
// For example: potato storage or client session storage
type Storage interface {
	// Create creates a record using the received data
	Create(interface{}) error
	// Obtain obtains record by id
	Obtain(string) (interface{}, error)
	// Delete removed record from the storage
	Delete(string) error
}

// StateStorage storage of states related to authorization requests
// Basically saves instances of model.Authorization
type StateStorage struct {
	*redis.Client
}

// generateKey creates a new key to be an identifier in the redis database based on the received state
func generateKey(state string) string {
	return "state:" + state
}

// Create receives a model.Authorization and uses your data to create a state-identified record
//
// Notes:
//
// - The record only is created if the state does not exist
//
// - If the record exists an error of type model.DuplicateRecord is returned
func (s StateStorage) Create(i interface{}) error {
	auth := i.(model.Authorization)

	cmd := s.SetNX(context.TODO(), generateKey(auth.State), model.BinaryJSON{I: auth}, 10*time.Second)

	flag, err := cmd.Result()
	if err != nil {
		return err
	}

	if !flag {
		err = model.DuplicateRecord(fmt.Sprintf(`a record to the state "%s" already exists`, auth.State))
	}

	return err
}

// Obtain search a saved instance of model.Authorization by the state
func (s StateStorage) Obtain(state string) (interface{}, error) {
	cmd := s.Get(context.TODO(), generateKey(state))

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
func (s StateStorage) Delete(state string) (err error) {
	return s.Del(context.TODO(), generateKey(state)).Err()
}

// ApplicationFinder defines a finder for application data
type ApplicationFinder interface {
	// FindApplication obtains a application by id
	FindApplication(string) (model.Application, error)
}
