package adapters

import (
	"github.com/equilibristofgo/sandbox/04_internal/app1/config"
)

type StandardService struct {
}

//ServiceHandler Contains the adapters services
type AdapterHandlerApp1 struct {
	Service1 StandardService
}

//NewServiceHandler instantiate the services of the adapters
func NewServiceHandler(cfg *config.App1Config) AdapterHandlerApp1 {
	return AdapterHandlerApp1{}
}
