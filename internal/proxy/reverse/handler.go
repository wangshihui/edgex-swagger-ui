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

package reverse

import (
	"fmt"
	bootsrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/gorilla/mux"
	"math"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
	"thundersoft.com/edgex/swagger-ui/internal/config"
	"thundersoft.com/edgex/swagger-ui/internal/container"
)

type proxyServers struct {
}

type EdgexReversProxy struct {
	dic          *di.Container
	proxyServers map[string][]*url.URL
	reversProxy  http.Handler
}

func NewEdgexReversProxy(dic *di.Container) *EdgexReversProxy {
	proxyServers := make(map[string][]*url.URL)

	config := container.ConfigurationFrom(dic.Get)

	ccn := len(config.Swagger.CoreComponents)

	if ccn > 0 {
		for _, c := range config.Swagger.CoreComponents {
			produceProxyServers(c, proxyServers, config)
		}
	}
	dcn := len(config.Swagger.DeviceComponents)

	if dcn > 0 {
		for _, c := range config.Swagger.DeviceComponents {
			produceProxyServers(c, proxyServers, config)
		}
	}

	return &EdgexReversProxy{
		dic:          dic,
		proxyServers: proxyServers,
		reversProxy:  reverseProxy(proxyServers, dic, config.Swagger.ProxyPrefix),
	}
}

func produceProxyServers(c config.ConfiComponent, proxyServers map[string][]*url.URL, config *config.ConfigurationStruct) {
	urls := make([]*url.URL, 0, math.MaxInt8)
	port := c.Port
	host := c.Host
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	if port != "" {
		host = host + ":" + port
	}
	urls = append(urls, &url.URL{
		Host:   host,
		Scheme: scheme,
	})
	proxyServers[path.Clean("/"+config.Swagger.ProxyPrefix+"/"+c.Route+"/")] = urls
}

func (p EdgexReversProxy) AddReversProxy(router *mux.Router) {
	for n, _ := range p.proxyServers {
		router.PathPrefix(n).Handler(p.reversProxy)
	}
}

func reverseProxy(ProxyServers map[string][]*url.URL, dic *di.Container, prefix string) *httputil.ReverseProxy {
	log := bootsrapContainer.LoggingClientFrom(dic.Get)
	director := func(req *http.Request) {
		if ProxyServers == nil {
			log.Info("no revers proxy configed")
		} else {
			reqPath := req.URL.Path
			for n, urls := range ProxyServers {
				if strings.Index(reqPath, n) == 0 {
					t := urls[rand.Int()%len(urls)]
					toPath := strings.Replace(reqPath, n, "", 1)
					req.URL.Scheme = t.Scheme
					req.URL.Host = t.Host
					req.URL.Path = path.Clean(toPath)
					log.Info(fmt.Sprintf("proxy request %s to %s", reqPath, req.URL))
					break
				}
			}
		}
	}
	return &httputil.ReverseProxy{Director: director}
}
