package main

import (
	"github.com/equilibristofgo/sandbox/03_mod_replace/pkg"

	"github.com/kataras/iris/v12"
)

var (
	pageSize                           = 4
	entityStorage pkg.EntityRepository = pkg.NewEntityMemoryRepository()
)

func main() {
	app := iris.New()
	entityStorage.Init()

	app.Logger().Fatal(app.Listen(":8080"))
}
