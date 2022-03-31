package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

// NewRedisClient establish connection to Redis using Configuration
func NewRedisClient(config Configuration) (client *redis.Client, err error) {
	db, err := strconv.Atoi(config.Database)
	if err != nil {
		return nil, fmt.Errorf("database must be a number")
	}

	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       db,
	})

	return client, client.Ping(context.TODO()).Err()
}
