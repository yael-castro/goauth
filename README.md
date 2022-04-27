# GOAuth - Authorization server
[![Icon](./doc/images/banner.png)](https://github.com/yael-castro)
[![Go Report Card](https://goreportcard.com/badge/github.com/yael-castro/goauth)](https://goreportcard.com/report/github.com/yael-castro/goauth)


Authorization server based on the [Authorization Code Flow](https://datatracker.ietf.org/doc/html/rfc6749#section-4.1) with the extension [Proof Key for Code Exchange (PKCE)](https://datatracker.ietf.org/doc/html/rfc7636) of the protocol [OAuth 2.0](https://datatracker.ietf.org/doc/html/rfc6749)


###### Optional features excluded
- Refresh tokens
- Redirect URL in the authorization response

<hr>

###### Architecture style explained
The architecture style used in this project is the most common layered architecture pattern
with a little changes.

```
internal
├── business    (business logic layer)
├── dependency  (manage dependencies)
├── handler     (presentation layer)
├── model       (data transfer objects, business objects, errors and enums)
└── repository  (persistence layer)
```
      
###### Required environment variables
[.env file example](./.env.example)

###### Configure your own private RSA key
```shell
export PRIVATE_RSA_KEY="$(openssl genrsa 1024)"
```

###### How to try
```go
package main

import (
	"github.com/yael-castro/goauth/internal/dependency"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.SetFlags(log.Flags() | log.Lshortfile)

	mux := http.NewServeMux()

	// Using dependency injection test profile you can test the server
	// without any configuration
	err := dependency.NewInjector(dependency.Testing).Inject(mux)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(`http server is running on port "%v" %v`, port, "🤘\n")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
```

<hr>
<a href="https://www.flaticon.com/free-icons/authentication" title="authentication icons">Authentication icons created by alkhalifi design - Flaticon</a>