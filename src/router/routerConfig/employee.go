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
		Authencication: true,
	},

	{
		URI:            "/employee",
		Method:         http.MethodGet,
		Function:       controllers.List,
		Authencication: true,
	},

	{
		URI:            "/employee/{employeeID}",
		Method:         http.MethodGet,
		Function:       controllers.ListID,
		Authencication: true,
	},

	{
		URI:            "/employee/{employeeID}",
		Method:         http.MethodPut,
		Function:       controllers.Update,
		Authencication: true,
	},

	{
		URI:            "/employee/{employeeID}",
		Method:         http.MethodDelete,
		Function:       controllers.Delete,
		Authencication: true,
	},

	{
		URI:            "/employee/{employeeID}/updatePassWord",
		Method:         http.MethodPost,
		Function:       controllers.UpdatePassWord,
		Authencication: true,
	},
}
