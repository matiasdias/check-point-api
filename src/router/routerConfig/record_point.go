package routerconfig

import (
	"check-point/src/controllers"
	"net/http"
)

var routesRecordPoint = []Routers{
	{
		URI:            "/recordPoint",
		Method:         http.MethodPost,
		Function:       controllers.Create,
		Authencication: false,
	},
}
