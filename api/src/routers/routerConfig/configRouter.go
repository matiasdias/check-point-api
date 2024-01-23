package routerconfig

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Routers struct {
	URI      string
	Method   string
	Function func(http.ResponseWriter, *http.Request)
}

// Config is used to place routes within the router
func Config(r *mux.Router) *mux.Router {
	return r
}
