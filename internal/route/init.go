package route

import (
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/gorilla/mux"
	"thundersoft.com/edgex/swagger-ui/internal/controllers"
)

func LoadRestRoutes( router *mux.Router, dic *di.Container)()  {
	sv := controllers.NewHttpResServer(dic)
	sv.Init()
	sv.AddHttpSwaggerHtmlHandler(router)
	sv.AddHttpSwaggerJsonHandler(router)
}
