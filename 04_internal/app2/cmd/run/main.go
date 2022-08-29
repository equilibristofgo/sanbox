package main

import (
	"github.com/equilibristofgo/sandbox/04_internal/app2/config"
	"github.com/equilibristofgo/sandbox/04_internal/app2/internal/adapters"
	"github.com/equilibristofgo/sandbox/04_internal/app2/internal/app"
	"github.com/equilibristofgo/sandbox/04_internal/app2/internal/ports"
)

//main Initialize each layer dependencies.
func main() {
	cfg := config.GetConfig()

	adaptersServiceHandler := adapters.NewServiceHandler(cfg)
	appServiceHandler := app.NewServiceHandler(adaptersServiceHandler.Service2)
	portsServiceHandler := ports.NewServiceHandler(cfg, appServiceHandler)
	portsServiceHandler.Cli.StartWorker()
}
