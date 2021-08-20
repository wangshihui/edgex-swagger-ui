package route

import (
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/gorilla/mux"
	"thundersoft.com/edgex/swagger-ui/internal/container"
	"thundersoft.com/edgex/swagger-ui/internal/controllers"
	"thundersoft.com/edgex/swagger-ui/internal/proxy/reverse"
)

func LoadRestRoutes(router *mux.Router, dic *di.Container) {
	sv := controllers.NewHttpResServer(dic)
	sv.Init()
	sv.AddHttpSwaggerHtmlHandler(router)
	sv.AddHttpSwaggerJsonHandler(router)

	config := container.ConfigurationFrom(dic.Get)
	if config.Swagger.Proxy {
		revers:=reverse.NewEdgexReversProxy(dic)
		revers.AddReversProxy(router)
	}

}
