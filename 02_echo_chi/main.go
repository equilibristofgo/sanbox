// Sample app to play with echo and chi, tryign to merge both frameworks
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"

	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"
)

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func main() {
	r := chi.NewRouter()
	r.Use(chiMW.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(echoMW.Logger())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

	http.ListenAndServe(":3000", r)
}
