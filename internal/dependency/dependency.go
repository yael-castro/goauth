// Package dependency manages
package dependency

import (
	"fmt"
	"github.com/yael-castro/godi/internal/business"
	"github.com/yael-castro/godi/internal/handler"
	"github.com/yael-castro/godi/internal/model"
	"github.com/yael-castro/godi/internal/repository"
	"net/http"
	"os"
	"strconv"
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

	generator := business.JWTGenerator{}

	err := generator.SetPrivateKey([]byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIBOgIBAAJBAKU2625UUvCglXhj9vzzHVJmlWKEAsYlQhK3+EvrVgpRDE7YQm4T
nYspJsYnVNQdNIdM0XRpSSe5p/pNDrjdx88CAwEAAQJAWpuCBqIMUpdfIgWA4Tzb
qeNEriDD/LNWRznJ3KkWKNVkhN8IQGItjDoPW5FzIbTplb47lj3Xa5/F7SNhh9l8
YQIhANOmy+f/cyl+o23omSUFq5amoNGhndC4nRbawLB2vUlnAiEAx9UxxG9DxlN+
neNDESR/o5XUO4HfPsJWbxlp4Dj1xVkCIFdL6biD5WUNBa10jY32m8JkcdplFamc
K7bcfTOLliErAiEAtlC75wvcOcVTb5k4RxuVmBnKV8BVfVywnwwAnKFbGYECIGIR
trcoqHkRGY43kVrARGwxQDv6+MGlWLCQ2m9p/mNv
-----END RSA PRIVATE KEY-----`))
	if err != nil {
		return err
	}

	grant := business.AuthorizationCodeGrant{
		TokenGenerator: generator,
		ScopeParser:    business.NewScopeParser(),
		CodeGenerator:  business.GenerateRandomCode,
		Owner: business.OwnerAuthenticator{
			Storage: &repository.MockStorage{
				"contacto@yael-castro.com": model.Owner{
					Id:       "contacto@yael-castro.com",
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
		CodeStorage:    &repository.MockStorage{},
		SessionStorage: &repository.MockStorage{},
		PKCE:           business.ProofKeyCodeExchange{},
	}

	*mux = *handler.NewServeMux(grant)
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

	generator := business.JWTGenerator{}

	err = generator.SetPrivateKey([]byte(os.Getenv("PRIVATE_RSA_KEY")))
	if err != nil {
		return err
	}

	redisClient, err := repository.NewRedisClient(redisSettings)
	if err != nil {
		return err
	}

	grant := &business.AuthorizationCodeGrant{
		TokenGenerator: generator,
		CodeGenerator:  business.CodeGeneratorFunc(business.GenerateUUID),
		SessionStorage: repository.SessionStorage{Client: redisClient},
		CodeStorage:    repository.StateStorage{Client: redisClient},
		PKCE:           business.ProofKeyCodeExchange{},
		Owner: business.OwnerAuthenticator{
			Storage: repository.OwnerStorage{Client: redisClient},
		},
		Client: business.ClientAuthenticator{
			Finder: repository.ClientFinder{Client: redisClient},
		},
		ScopeParser: business.NewScopeParser(),
	}

	*mux = *handler.NewServeMux(grant)
	return nil
}
