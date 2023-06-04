package main

import (
	"github.com/equilibristofgo/sandbox/04_internal/app1/config"
	"github.com/equilibristofgo/sandbox/04_internal/app1/internal/adapters"
	"github.com/equilibristofgo/sandbox/04_internal/app1/internal/app"
	"github.com/equilibristofgo/sandbox/04_internal/app1/internal/ports"
)

//main Initialize each layer dependencies.
func main() {
	cfg := config.GetConfig()

	adaptersServiceHandler := adapters.NewServiceHandler(cfg)
	appServiceHandler := app.NewServiceHandler(adaptersServiceHandler.Service1)
	portsServiceHandler := ports.NewServiceHandler(cfg, appServiceHandler)
	portsServiceHandler.Cli.StartWorker()
}
