package djijane

import (
	"dji-joe"

	"net/http"
	"time"
	
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func HttpLogger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		Log.InfoF("%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = HttpLogger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{ 
	Route{
		"Heartbeat",
		"POST",
		djijoe.API_HEARTBEAT,
		Heartbeat,
	},

	Route{
		"WakeUp",
		"POST",
		djijoe.API_WAKEUP,
		WakeUp,
	},

	Route{
		"ShutDown",
		"POST",
		djijoe.API_SHUTDOWN,
		ShutDown,
	},

	Route{
		"NewDroneInfo",
		"POST",
		djijoe.API_NEWDRONEINFO,
		NewDroneInfo,
	},
}
