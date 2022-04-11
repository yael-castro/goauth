// Package dependency manages
package dependency

import (
	"fmt"
	"github.com/yael-castro/godi/internal/business"
	"github.com/yael-castro/godi/internal/handler"
	"github.com/yael-castro/godi/internal/model"
	"github.com/yael-castro/godi/internal/repository"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Profile defines options of dependency injection
type Profile uint

// Supported profiles for dependency injection
const (
	// Default defines the production profile
	Default Profile = iota
	Testing
)

// Injector defines a dependency injector
type Injector interface {
	// Inject takes any data type and fill of required dependencies (dependency injection)
	Inject(interface{}) error
}

// InjectorFunc function that implements the Injector interface
type InjectorFunc func(interface{}) error

func (f InjectorFunc) Inject(i interface{}) error {
	return f(i)
}

// NewInjector is an abstract factory to Injector, it builds a instance of Injector interface based on the Profile based as parameter
//
// Supported profiles: Default and Testing
//
// If pass a parameter an invalid profile it panics
func NewInjector(p Profile) Injector {
	switch p {
	case Testing:
		return InjectorFunc(testingProfile)
	case Default:
		return InjectorFunc(defaultProfile)
	}

	panic(fmt.Sprintf(`invalid profile: "%d" is not supported`, p))
}

// testingProfile InjectorFunc for *handler.Handler that uses a Testing Profile
func testingProfile(i interface{}) error {
	mux, ok := i.(*http.ServeMux)
	if !ok {
		return fmt.Errorf(`invalid type "%T"`, i)
	}

	random := rand.New(rand.NewSource(time.Now().Unix()))

	authorizer := business.OAuth{
		"code": business.AuthorizationCodeGrant{
			CodeGenerator: business.CodeGeneratorFunc(func() model.AuthorizationCode {
				return model.AuthorizationCode(strconv.Itoa(random.Int()))
			}),
			Owner: business.OwnerAuthenticator{
				Storage: &repository.MockStorage{
					"contacto@yael-castro.com": model.Owner{
						Password: "$2a$10$g141w.TTnp5Bm/rLNqRRRevOSFhKBdV5KaJYxEDi9U5R9TgkZbfne", // yael.castro
					},
				},
			},
			Client: business.ClientAuthenticator{
				Finder: repository.MockClientFinder{
					"mobile": model.Client{
						Id: "mobile",
						AllowedOrigins: []string{
							"http://localhost/callback",
							"http://localhost:8080/callback",
						},
					},
				},
			},
			Storage: &repository.MockStorage{},
			PKCE:    business.ProofKeyCodeExchange{},
		},
	}

	mux.HandleFunc("/go-auth/v1/authorization", handler.NewAuthorizationHandler(authorizer))

	return nil
}

// defaultProfile InjectorFunc for *handler.Handler that uses a Testing Profile
func defaultProfile(i interface{}) error {
	mux, ok := i.(*http.ServeMux)
	if !ok {
		return fmt.Errorf(`invalid type "%T"`, i)
	}

	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		return err
	}

	redisSettings := repository.Configuration{
		Type:     repository.KeyValue,
		Host:     os.Getenv("REDIS_HOST"),
		Port:     redisPort,
		Database: os.Getenv("REDIS_DATABASE"),
		User:     os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	redisClient, err := repository.NewRedisClient(redisSettings)
	if err != nil {
		return err
	}

	authorizer := business.OAuth{
		"code": business.AuthorizationCodeGrant{
			CodeGenerator: business.CodeGeneratorFunc(business.NewUUIDCode),
			Storage:       repository.StateStorage{Client: redisClient},
			PKCE:          business.ProofKeyCodeExchange{},
			Owner: business.OwnerAuthenticator{
				Storage: repository.OwnerStorage{Client: redisClient},
			},
			Client: business.ClientAuthenticator{
				Finder: repository.ClientFinder{Client: redisClient},
			},
		},
	}

	mux.HandleFunc("/go-auth/v1/authorization", handler.NewAuthorizationHandler(authorizer))

	return nil
}
