package routerconfig

import (
	"check-point/src/controllers"
	"net/http"
)

var routesEmployee = []Routers{
	//CRUD DE FUNCIONARIOS
	{
		URI:            "/employee",
		Method:         http.MethodPost,
		Function:       controllers.Create,
		Authencication: false,
	},

	{
		URI:            "/employee",
		Method:         http.MethodGet,
		Function:       controllers.List,
		Authencication: false,
	},
}
