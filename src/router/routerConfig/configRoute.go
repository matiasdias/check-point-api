package routerconfig

import (
	"check-point/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

// Routers representa todas as rotas da api
type Routers struct {
	URI            string
	Method         string
	Function       func(http.ResponseWriter, *http.Request)
	Authencication bool
}

// Config serve para colocando as rotas dentro do router
func Config(r *mux.Router) *mux.Router {
	routers := routesEmployee
	routers = append(routers, routeLogin)
	routers = append(routers, routesRecordPoint...)

	for _, router := range routers {

		if router.Authencication {
			r.HandleFunc(router.URI, middlewares.Logger(middlewares.Authenticate(router.Function))).
				Methods(router.Method)
		} else {
			r.HandleFunc(router.URI, middlewares.Logger(router.Function)).Methods(router.Method)
		}
	}

	return r
}
