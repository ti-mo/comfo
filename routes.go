package main

import (
	"github.com/gorilla/mux"
	"net/http"
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

var routes = Routes{
	Route{
		"GetTemperatures",
		"GET",
		"/temps",
		TempHandlerGet,
	},
	Route{
		"GetFans",
		"GET",
		"/fans",
		FanHandlerGet,
	},
	Route{
		"SetFans",
		"PUT",
		"/fans/{speed:(?:[1-4]|up|down)}",
		FanHandlerSet,
	},
	Route{
		"GetFanProfiles",
		"GET",
		"/profiles",
		FanProfilesHandlerGet,
	},
	Route{
		"FlushCache",
		"PUT",
		"/cache/flush/{cache:[a-z]+}",
		FlushCacheHandler,
	},
}
