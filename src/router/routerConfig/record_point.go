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
		Authencication: false,
	},

	{
		URI:            "/recordPoint",
		Method:         http.MethodGet,
		Function:       controllers.ListRecordPoint,
		Authencication: false,
	},
}
