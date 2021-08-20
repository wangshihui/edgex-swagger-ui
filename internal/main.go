package internal

import (
	"context"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/flags"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/handlers"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/interfaces"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/gorilla/mux"
	"os"
	edgex "thundersoft.com/edgex/swagger-ui"
	"thundersoft.com/edgex/swagger-ui/internal/config"
	"thundersoft.com/edgex/swagger-ui/internal/container"
)

func Main(ctx context.Context, cancel context.CancelFunc, router *mux.Router) {
	startupTimer := startup.NewStartUpTimer(common.CoreMetaDataServiceKey)

	// All common command-line flags have been moved to DefaultCommonFlags. Service specific flags can be add here,
	// by inserting service specific flag prior to call to commonFlags.Parse().
	// Example:
	// 		flags.FlagSet.StringVar(&myvar, "m", "", "Specify a ....")
	//      ....
	//      flags.Parse(os.Args[1:])
	//
	f := flags.New()
	f.Parse(os.Args[1:])

	configuration := &config.ConfigurationStruct{}
	dic := di.NewContainer(di.ServiceConstructorMap{
		container.ConfigurationName: func(get di.Get) interface{} {
			return configuration
		},
	})

	httpServer := handlers.NewHttpServer(router, true)
	bootstrap.Run(
		ctx,
		cancel,
		f,
		common.CoreMetaDataServiceKey,
		ConfigStemCore,
		configuration,
		startupTimer,
		dic,
		false,
		[]interfaces.BootstrapHandler{
			NewBootstrap(router).BootstrapHandler,
			httpServer.BootstrapHandler,
			handlers.NewStartMessage("edgex-swagger-ui-security", edgex.Version).BootstrapHandler,
		})
}
