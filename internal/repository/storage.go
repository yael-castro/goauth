package repository

import "github.com/yael-castro/godi/internal/model"

// SessionStorage defines a storage for any session type
type SessionStorage interface {
	// CreateSession creates a session based on the received data
	CreateSession(interface{}) error
	// ObtainSession obtains session by id
	ObtainSession(string) (interface{}, error)
	// DeleteSession revokes a session
	DeleteSession(string) error
}

// ApplicationFinder defines a finder for application data
type ApplicationFinder interface {
	// FindApplication obtains a application by id
	FindApplication(string) (model.Application, error)
}
