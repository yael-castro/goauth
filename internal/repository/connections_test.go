package repository

import (
	"strconv"
	"testing"
)

// TestNewRedisClient health check for redis connections
func TestNewRedisClient(t *testing.T) {
	tdt := []struct {
		config      Configuration
		expectedErr error
	}{
		// Test case for default connection
		{
			config: Configuration{
				Type:     KeyValue,
				Host:     "localhost",
				Port:     6379,
				Database: "0",
			},
		},
	}

	for i, v := range tdt {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			_, err := NewRedisClient(v.config)
			if err != v.expectedErr {
				t.Fatal(err)
			}

			t.Log(v.config)
		})
	}

}
