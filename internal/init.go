package internal

import (
	"context"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/gorilla/mux"
	"sync"
	"thundersoft.com/edgex/swagger-ui/internal/route"
)

// Bootstrap  contains references to dependencies required by the BootstrapHandler.
// template code from edgex core
type Bootstrap struct {
	router     *mux.Router
}

// NewBootstrap is a factory method that returns an initialized Bootstrap receiver struct.
func NewBootstrap(router *mux.Router ) *Bootstrap {
	return &Bootstrap{
		router:     router,
	}
}

// BootstrapHandler fulfills the BootstrapHandler contract and performs initialization needed by the command service.
func (b *Bootstrap) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup, _ startup.Timer, dic *di.Container) bool {
	route.LoadRestRoutes(b.router, dic)
	return true
}
