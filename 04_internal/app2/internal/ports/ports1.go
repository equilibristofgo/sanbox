package ports

import (
	"github.com/equilibristofgo/sandbox/04_internal/app2/config"
	"github.com/equilibristofgo/sandbox/04_internal/app2/internal/app"
)

type Foo struct {
	Bar func()
}

type StandardClient struct {
	StartWorker func()
}

//ServiceHandler Contains the ports services
type ServiceHandler struct {
	cfg *config.App2Config
	Cli *StandardClient
}

//NewServiceHandler instantiate the services of the ports
func NewServiceHandler(cfg *config.App2Config, handler *app.AppHandlerApp2) ServiceHandler {
	return ServiceHandler{cfg: cfg, Cli: nil}
}
