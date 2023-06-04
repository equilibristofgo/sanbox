package app

import (
	"github.com/equilibristofgo/sandbox/04_internal/app1/internal/adapters"
)

type AppHandlerApp1 struct {
	service1 adapters.StandardService
}

//NewServiceHandler Initialize the application layer services
func NewServiceHandler(service1 adapters.StandardService) *AppHandlerApp1 {
	return &AppHandlerApp1{
		service1,
	}

}
