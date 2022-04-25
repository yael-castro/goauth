package main

import (
	"github.com/yael-castro/go-auth/internal/dependency"
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

	err := dependency.NewInjector(dependency.Testing).Inject(mux)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(`http server is running on port "%v" %v`, port, "🤘\n")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
