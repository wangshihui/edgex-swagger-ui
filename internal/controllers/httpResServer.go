package controllers

import (
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"net/http"
)

type HttpResServer struct {
	dic *di.Container

}

func NewHttpResServer(dic *di.Container) HttpResServer {
	return HttpResServer{
		dic:dic,
	}
}


// GetDeviceByType
func (dc *HttpResServer) GetDeviceDefByName(writer http.ResponseWriter, request *http.Request) {
	// URL parameters
	//vars := mux.Vars(request)
	//deviceType := vars[common.Name]
	log:=container.LoggingClientFrom(dc.dic.Get)
	log.Info(request.RequestURI)

}