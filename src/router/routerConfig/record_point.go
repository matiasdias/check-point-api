package routerconfig

import (
	"check-point/src/controllers"
	"net/http"
)

var routesRecordPoint = []Routers{
	{
		URI:            "/recordPoint",
		Method:         http.MethodPost,
		Function:       controllers.CreateRecordPoint,
		Authencication: true,
	},

	{
		URI:            "/recordPoint",
		Method:         http.MethodGet,
		Function:       controllers.ListRecordPoint,
		Authencication: true,
	},

	{
		URI:            "/recordPoint/{recordEmployeeID}",
		Method:         http.MethodGet,
		Function:       controllers.ListIDRecordPoint,
		Authencication: true,
	},

	{
		URI:            "/recordPoint/{recordEmployeeID}",
		Method:         http.MethodPut,
		Function:       controllers.UpdateRecordPoint,
		Authencication: true,
	},

	{
		URI:            "/recordPoint/{recordEmployeeID}",
		Method:         http.MethodDelete,
		Function:       controllers.DeleteRecordPoint,
		Authencication: true,
	},
}
