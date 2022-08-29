package ports

import (
	"fmt"

	"github.com/equilibristofgo/sandbox/04_internal/app1/config"
	"github.com/equilibristofgo/sandbox/04_internal/app1/internal/app"
)

type Foo struct {
	Bar func()
}

type StandardClient struct {
	StartWorker func()
}

//ServiceHandler Contains the ports services
type ServiceHandler struct {
	cfg *config.App1Config
	Cli *StandardClient
}

//NewServiceHandler instantiate the services of the ports
func NewServiceHandler(cfg *config.App1Config, handler *app.AppHandlerApp1) ServiceHandler {
	return ServiceHandler{cfg: cfg, Cli: &StandardClient{
		func() { fmt.Println("remote app1") },
	}}
}
