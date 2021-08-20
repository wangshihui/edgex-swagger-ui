/*******************************************************************************
 * Copyright wangshihui
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/


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
	if config.Swagger.ReverseProxy {
		revers := reverse.NewEdgexReversProxy(dic)
		revers.AddReversProxy(router)
	}

}
