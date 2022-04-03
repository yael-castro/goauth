package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/yael-castro/godi/internal/model"
)

// Finder defines a finder for saved data in some storage
type Finder interface {
	// Find obtains a record by id
	Find(string) (interface{}, error)
}

// _ "implement" constraint for ClientFinder
var _ Finder = ClientFinder{}

// ClientFinder creates a Finder implementation to find OAuth clients
type ClientFinder struct {
	*redis.Client
}

// clientKey creates a key with the pattern "client:<clientId>" to save a client
func (ClientFinder) clientKey(clientId string) string {
	return "client:" + clientId
}

// secretKey creates a key with the pattern "client:<clientId>:secret" to save a client secret
func (c ClientFinder) secretKey(clientId string) string {
	return c.clientKey(clientId) + ":secret"
}

// listKey creates a key with the pattern "client:<clientId>:origins" to save allowed origins for client
func (c ClientFinder) listKey(clientId string) string {
	return c.clientKey(clientId) + ":origins"
}

// Find search a client by client id
func (c ClientFinder) Find(clientId string) (i interface{}, err error) {
	client := model.Client{Id: clientId}

	client.Secret, err = c.Get(context.TODO(), c.secretKey(clientId)).Result()
	if err != nil {
		return
	}

	client.AllowedOrigins, err = c.LRange(context.TODO(), c.listKey(clientId), 0, 10).Result()
	i = client
	return
}

// _ "implement" constraint for MockClientFinder
var _ Finder = MockClientFinder{}

// MockClientFinder mock store for model.Client
type MockClientFinder map[string]model.Client

// Find search a client by id
func (m MockClientFinder) Find(clientId string) (interface{}, error) {
	client, ok := m[clientId]
	if !ok {
		return nil, model.NotFound(fmt.Sprintf(`missing client "%s"`, clientId))
	}

	client.Id = clientId

	return client, nil
}
