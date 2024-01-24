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
}
