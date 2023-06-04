package app

import (
	"github.com/equilibristofgo/sandbox/04_internal/app2/internal/adapters"
)

type AppHandlerApp2 struct {
	service2 adapters.StandardService
}

//NewServiceHandler Initialize the application layer services
func NewServiceHandler(service2 adapters.StandardService) *AppHandlerApp2 {
	return &AppHandlerApp2{
		service2,
	}

}
