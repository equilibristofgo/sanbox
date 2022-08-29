package main

import (
	exporter "github.com/equilibristofgo/sandbox/04_internal/app1/cmd"
	"github.com/equilibristofgo/sandbox/04_internal/app1/config"
)

//main Initialize each layer dependencies.
func main() {
	_ = config.GetConfig()

	e := exporter.Exporter{}
	e.Exported()
}
