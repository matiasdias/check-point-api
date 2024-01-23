package routers

import (
	routerconfig "check-point/api/src/routers/routerConfig"

	"github.com/gorilla/mux"
)

func Load() *mux.Router {
	r := mux.NewRouter()
	return routerconfig.Config(r)
}
