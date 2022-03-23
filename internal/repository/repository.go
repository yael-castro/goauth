package repository

import "fmt"

// Type defines the repository type
type Type uint

// Supported values for Type
const (
	SQL Type = iota
	// NoSQL defines an NoSQL storage oriented for documents, more specifically MongoDB
	NoSQL
	KeyValue
)

// _ "implement" constraint for Configuration struct
var _ fmt.Stringer = Configuration{}

// Configuration options used to establish conenction to any repository
type Configuration struct {
	// Type defines the type of repository to which you want to establish a connection
	Type
	Host     string
	Port     int
	Database string
	User     string
	// Password is a password of repository user but also could be a token to authenticate
	Password string
	// Secure defines options like +srv in MongoDB or SSL in Postgresql
	Secure bool
	// Debug defines if the debug mode of repository operations are enabled
	Debug bool
}

// String build and returns a URI used to establish a connection to any repository defined by the Type embbed in the Configuration structure
//
// Supported types KeyValue and NoSQL
//
// KeyValue returns "redis://<c.User>:<c.Password>@<c.Host>:<c.Port>/<c.Database>"
//
// NoSQL returns "mongodb+srv://<c.User>:<c.Password>@<c.Host>" or "mongodb://<c.User>:<c.Password>@<c.Host>:<c.Port>"
//
func (c Configuration) String() string {
	switch c.Type {
	case NoSQL:
		if c.Secure {
			return fmt.Sprintf(
				"mongodb+srv://%s:%s@%s", //?maxPoolSize=%s",
				c.User,
				c.Password,
				c.Host,
			)
		}

		return fmt.Sprintf(
			"mongodb://%s:%s@%s:%d", //?maxPoolSize=%s",
			c.User,
			c.Password,
			c.Host,
			c.Port,
		)
	case KeyValue:
		return fmt.Sprintf("redis://%s:%s@%s:%d/%s", c.User, c.Password, c.Host, c.Port, c.Database)
	}

	panic(fmt.Sprintf(`type "%d" is not supported`, c.Type))
}
