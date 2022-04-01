package repository

import (
	"context"
	"encoding/json"
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

// clientKey creates a key to search a client
func (ClientFinder) clientKey(clientId string) string {
	return "client:" + clientId
}

// listKey creates a key to search the allowed redirect url of client
func (ClientFinder) listKey(clientId string) string {
	return "client:" + clientId + ":origins"
}

// Find search a client by client id
func (r ClientFinder) Find(clientId string) (i interface{}, err error) {
	client := model.Client{}

	serializedData, err := r.Get(context.TODO(), r.clientKey(clientId)).Result()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(serializedData), &client)
	if err != nil {
		return
	}

	client.AllowedOrigins, err = r.LRange(context.TODO(), r.listKey(clientId), 0, 10).Result()
	i = client
	return
}
