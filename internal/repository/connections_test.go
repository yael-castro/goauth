package repository

import (
	"strconv"
	"testing"
)

var defaultRedisConfiguration = Configuration{
	Type:     KeyValue,
	Host:     "localhost",
	Port:     6379,
	Database: "0",
}

// TestNewRedisClient health check for redis connections
func TestNewRedisClient(t *testing.T) {
	tdt := []struct {
		config      Configuration
		expectedErr error
	}{
		// Test case for default connection
		{config: defaultRedisConfiguration},
	}

	for i, v := range tdt {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			client, err := NewRedisClient(v.config)
			if err != v.expectedErr {
				t.Fatal(err)
			}

			t.Cleanup(func() {
				_ = client.Close()
			})

			if client == nil {
				t.Fatal("redis client returned is nil")
			}

			t.Log(v.config)
		})
	}
}
