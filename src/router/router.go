package router

import (
	routerconfig "check-point/src/router/routerConfig"

	"github.com/gorilla/mux"
)

// Load vai retornar um router com as rotas configuradas
func Load() *mux.Router {
	r := mux.NewRouter()
	return routerconfig.Config(r)
}
