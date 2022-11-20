package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aaletov/go-smo/api"
	"github.com/aaletov/go-smo/server"
	"github.com/go-chi/chi/v5"
)

const (
	sourcesLambda = 13
	sourcesCount  = 3
	bufferCount   = 4
	devicesCount  = 3
)

func main() {
	var port = 8081

	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create an instance of our handler which satisfies the generated interface
	smoServer := server.NewServer()

	// This is how you set up a basic chi router
	r := chi.NewRouter()

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	//r.Use(middleware.OapiRequestValidator(swagger))

	// We now register our petStore above as the handler for the interface
	api.HandlerFromMux(smoServer, r)

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
