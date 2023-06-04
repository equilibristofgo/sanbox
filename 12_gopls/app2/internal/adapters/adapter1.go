package adapters

import (
	"github.com/equilibristofgo/sandbox/04_internal/app2/config"
)

type StandardService struct {
}

//ServiceHandler Contains the adapters services
type AdapterHandlerApp2 struct {
	Service2 StandardService
}

//NewServiceHandler instantiate the services of the adapters
func NewServiceHandler(cfg *config.App2Config) AdapterHandlerApp2 {
	return AdapterHandlerApp2{}
}
