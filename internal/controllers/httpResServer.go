package controllers

import (
	"encoding/json"
	"fmt"
	bootsrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"thundersoft.com/edgex/swagger-ui/internal/config"
	"thundersoft.com/edgex/swagger-ui/internal/container"
)

const (
	place_holder = `<<INSERT-SWAGGER-URLS>>`
	js_template  = `window.onload = function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
        urls: ` + place_holder + `,
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      });
      // End Swagger UI call region

      window.ui = ui;
    };`

	SwaggerDataRequest = "swagger"
	SwaggerJsFileName  = "edgex-swagger-init.js"

	swaggerComponetProperty = "components"
	securitySchemesProperty = "securitySchemes"
	swaggersecurityProperty = "security"

	defaultSchemaName = "Authorization"

	edgexKongAuth = "ApiKeyAuth"

	EdgexProxyPrefix = "ep"
)

type ApiKeyAuth struct {
	Type string `json:"type" :"type"`
	In   string `json:"in" :"type"`
	Name string `json:"name" :"type"`
}

type swaggerUrl struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type HttpResServer struct {
	dic                *di.Container
	swaggerFileHandler http.Handler
	prefix             string
	OpenApiJsonData    map[string]interface{}
	init               bool
}

func NewHttpResServer(dic *di.Container) HttpResServer {
	config := container.ConfigurationFrom(dic.Get)
	swagger := config.Swagger
	prefix := swagger.SwaggerPathPrefix
	if prefix == "" {
		prefix = "/"
	}
	swaggerDir := swagger.SwaggerFileDir
	if swaggerDir == "" {
		panic("swagger static file dir is empty")
	}

	if _, error := pathExists(swaggerDir); error != nil {
		panic("swagger static file dir not exist")
	}
	h := http.StripPrefix(prefix, http.FileServer(http.Dir(swaggerDir)))
	return HttpResServer{
		dic:                dic,
		swaggerFileHandler: h,
		prefix:             prefix,
	}
}

func (hs *HttpResServer) AddHttpSwaggerHtmlHandler(router *mux.Router) {
	router.PathPrefix(hs.prefix).HandlerFunc(hs.httpSwaggerFileHandler)
}
func (hs *HttpResServer) AddHttpSwaggerJsonHandler(router *mux.Router) {
	router.HandleFunc("/"+SwaggerDataRequest+"/"+"{"+common.Name+"}", hs.getSwaggerJsonData)
}

//httpSwaggerFileHandler
func (hs *HttpResServer) httpSwaggerFileHandler(writer http.ResponseWriter, request *http.Request) {
	hs.swaggerFileHandler.ServeHTTP(writer, request)
}

//getSwaggerJsonData
func (hs *HttpResServer) getSwaggerJsonData(writer http.ResponseWriter, request *http.Request) {
	log := bootsrapContainer.LoggingClientFrom(hs.dic.Get)
	vars := mux.Vars(request)
	componetName := vars[common.Name]
	log.Info(request.RequestURI, componetName)
	data, ok := hs.OpenApiJsonData[componetName]
	if ok {
		sendResponse(writer, request, "getSwaggerJsonData", log, data, http.StatusOK)
	} else {
		sendResponse(writer, request, "getSwaggerJsonData", log, "", http.StatusNotFound)
	}
}

func (hs *HttpResServer) Init() {
	log := bootsrapContainer.LoggingClientFrom(hs.dic.Get)
	if hs.init {
		log.Info(fmt.Sprintf("swagger swaggerServer has been inited"))
		return
	}

	config := container.ConfigurationFrom(hs.dic.Get)
	swagger := config.Swagger
	serviceHost := "//" + config.Service.Host + ":" + strconv.Itoa(config.Service.Port)
	swaggerServer := serviceHost

	if config.Swagger.ProxyPrefix == "" {
		config.Swagger.ProxyPrefix = EdgexProxyPrefix
	}

	if swagger.ReverseProxy {
		swaggerServer = swaggerServer + path.Clean("/"+swagger.ProxyPrefix+"/")
		log.Info(fmt.Sprintf("use proxy mode %s", swaggerServer))
	} else {
		swaggerServer = "//" + config.KongURL.Server + ":" + strconv.Itoa(config.KongURL.ApplicationPort)
	}
	log.Info(fmt.Sprintf("make swagger use base url '%s' proxy mod '%v'", swaggerServer, swagger.ReverseProxy))
	// then load yamls
	coreDir := swagger.CoreDir
	dvDir := swagger.DeviceSdkDir
	log.Info(fmt.Sprintf("load swagger yamls from dictionary [core]  %s  [deviceSerive template] %s ", coreDir, dvDir))

	if ok, _ := pathExists(coreDir); !ok {
		panic("core dir is not exist")
	}
	if ok, _ := pathExists(dvDir); !ok {
		panic("device service dir is not exist")
	}
	hs.OpenApiJsonData = make(map[string]interface{})
	swaggerUrls := make([]swaggerUrl, 0, len(swagger.CoreComponents)+len(swagger.DeviceComponents))
	for _, c := range swagger.CoreComponents {
		e, s, m := loadYamls(c, coreDir)
		if e != nil {
			log.Error(fmt.Sprintf("err when load yaml for componets e:= %s, c:= %s", e, c.Name))
			continue
		}
		swaggerUrls = append(swaggerUrls, swaggerUrl{
			Url:  serviceHost + path.Clean("/"+SwaggerDataRequest+"/"+c.Name),
			Name: c.Name,
		})
		processJson(m, s, c, swaggerServer, swagger.ReverseProxy)
		hs.OpenApiJsonData[c.Name] = m
	}
	for _, c := range swagger.DeviceComponents {
		e, s, m := loadYamls(c, dvDir)
		if e != nil {
			log.Error(fmt.Sprintf("err when load yaml for componets e:= %s, c:= s%", e, c.Name))
			continue
		}
		swaggerUrls = append(swaggerUrls, swaggerUrl{
			Url:  serviceHost + path.Clean("/"+SwaggerDataRequest+"/"+c.Name),
			Name: c.Name,
		})
		processJson(m, s, c, swaggerServer, swagger.ReverseProxy)
		hs.OpenApiJsonData[c.Name] = m
	}
	e := genInitJs(swaggerUrls, log, swagger.SwaggerFileDir)
	if e != nil {
		log.Errorf(fmt.Sprintf("error happen when gen edgex swagger javascript %s", e))
		panic("edgex swagger javascript gen error ")
	}
	hs.init = true
}

func processJson(m map[string]interface{}, s string, c config.ConfiComponent, server string, proxy bool) {
	//return
	var swaggerComponet map[string]interface{}

	_, ok := m[swaggerComponetProperty]
	if _, o := m[swaggerComponetProperty].(map[string]interface{}); !o || !ok {
		swaggerComponet = make(map[string]interface{})
		m[swaggerComponetProperty] = swaggerComponet
	}
	swaggerComponet, _ = m[swaggerComponetProperty].(map[string]interface{})
	_, ok = swaggerComponet[securitySchemesProperty]
	if !ok {
		swaggerComponet[securitySchemesProperty] = make(map[string]interface{})
	}
	sc, _ := swaggerComponet[securitySchemesProperty]
	t, o := sc.(map[string]interface{})
	if o {
		t[edgexKongAuth] = ApiKeyAuth{
			Type: "apiKey",
			In:   "header",
			Name: defaultSchemaName,
		}
	}

	secs := make([]interface{}, 0, 10)
	api := make(map[string]interface{})
	api[edgexKongAuth] = make([]string, 0, 1)
	secs = append(secs, api)
	m[swaggersecurityProperty] = secs
	//m[swaggersecurityProperty]=spv
	servers := make([]interface{}, 0, 1)

	apiServer := make(map[string]string)
	apiServer["url"] = server + path.Clean("/"+c.Route+"/"+c.ApiVer+"/")
	apiServer["description"] = " Use Kong GateWay"
	if proxy {
		apiServer["description"] = "Use Local ReversProxy"
	}
	servers = append(servers, apiServer)
	m["servers"] = servers
}

func genInitJs(urls []swaggerUrl, l logger.LoggingClient, swaggerDir string) error {
	us, e := json.Marshal(urls)
	if e != nil {
		l.Errorf(fmt.Sprintf("error happen when serialize url arry %s", e))
		return e
	}
	l.Info(fmt.Sprintf("swagger ui urls %s", us))
	wf := swaggerDir + string(os.PathSeparator) + SwaggerJsFileName
	if _, e := os.Stat(wf); e == nil {
		e = os.Remove(wf)
		if e != nil {
			l.Errorf(fmt.Sprintf("can note remove exist file %s %s", wf, e))
			return e
		}
	}

	js := strings.Replace(js_template, place_holder, string(us), 1)
	e = os.WriteFile(wf, []byte(js), fs.ModePerm)
	return e
}

func loadYamls(c config.ConfiComponent, dir string) (error, string, map[string]interface{}) {
	//files,err:=os.ReadDir(dir)
	//if err!=nil{
	//	return err,nil
	//}
	fp := dir + string(os.PathSeparator) + c.FileName
	result := make(map[string]interface{})
	_, e := os.Stat(fp)
	if e != nil {
		return e, "", nil
	}
	b, e := os.ReadFile(fp)
	if e != nil {
		return e, "", nil
	}
	yaml.Unmarshal(b, result)

	server := ""
	if c.Host != "" {
		server = c.Host + c.Port
	}

	return nil, server, result
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// sendResponse puts together the response packet for the V2 API
func sendResponse(
	writer http.ResponseWriter,
	request *http.Request,
	api string,
	lc logger.LoggingClient,
	response interface{},
	statusCode int) {

	correlationID := request.Header.Get(common.CorrelationHeader)

	writer.Header().Set(common.CorrelationHeader, correlationID)
	writer.Header().Set(common.ContentType, common.ContentTypeJSON)
	writer.WriteHeader(statusCode)

	data, err := json.Marshal(response)
	if err != nil {
		lc.Error(fmt.Sprintf("Unable to marshal %s response", api), "error", err.Error(), common.CorrelationHeader, correlationID)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(data)
	if err != nil {
		lc.Error(fmt.Sprintf("Unable to write %s response", api), "error", err.Error(), common.CorrelationHeader, correlationID)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
