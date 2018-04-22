package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ti-mo/comfo/comfoserver"
	rpc "github.com/ti-mo/comfo/rpc/comfo"
)

// Route is a gorilla/mux route entry.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// NewRouter is the initialization method for a gorilla/mux router.
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	// Add Twirp handler
	server := &comfoserver.Server{}
	twirpHandler := rpc.NewComfoServer(server, nil)
	router.PathPrefix(rpc.ComfoPathPrefix).Handler(Logger(twirpHandler, "twirp"))

	// Add Route entries to router
	for _, route := range routes {

		// Add middleware handlers
		handler := Logger(route.HandlerFunc, route.Name)
		handler = Headers(handler)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

// Routes is a list of Route items.
type Routes []Route

var routes = Routes{}
