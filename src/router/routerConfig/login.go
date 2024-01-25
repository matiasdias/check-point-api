package routerconfig

import (
	"check-point/src/controllers"
	"net/http"
)

var routeLogin = Routers{

	URI:            "/auth",
	Method:         http.MethodPost,
	Function:       controllers.Login,
	Authencication: false,
}
