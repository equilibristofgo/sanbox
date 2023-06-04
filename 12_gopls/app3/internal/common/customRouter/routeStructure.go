package customRouter

import (
	"github.com/go-chi/chi"
	"net/http"
)

//RouteStructure -
type RouteStructure struct {
	Uses    []func(http.Handler) http.Handler
	Methods []MethodRoutePatternInfo
	Mounts  []RoutePatternInfo
	Routes  map[string]*RouteStructure
	Groups  []RouteStructure
}

//MethodRoutePatternInfo -
type MethodRoutePatternInfo struct {
	Method  string // It's recommend to use constants of net/httpHandler/method.go
	Pattern string
	Handler http.Handler
}

//RoutePatternInfo -
type RoutePatternInfo struct {
	Pattern     string
	OtherParams []interface{}
	ServeFunc   func(router *chi.Mux, otherParams ...interface{})
}

//PrepareRoutes -
func (routeStructure RouteStructure) PrepareRoutes(router *chi.Mux) {
	router.Use(routeStructure.Uses...)

	for _, method := range routeStructure.Methods {
		//fmt.Println("Method:", method)
		router.Method(method.Method, method.Pattern, method.Handler)
	}

	for _, mountInfo := range routeStructure.Mounts {
		newRouter := chi.NewRouter()
		mountInfo.ServeFunc(newRouter, mountInfo.OtherParams...)
		router.Mount(mountInfo.Pattern, newRouter)
	}
	for pattern, routeInfo := range routeStructure.Routes {
		router.Route(pattern, func(newRouter chi.Router) {
			(*routeInfo).prepareRoutesNewRouter(newRouter)
		})
	}

	for _, groupInfo := range routeStructure.Groups {
		router.Group(func(newRouter chi.Router) {
			groupInfo.prepareRoutesNewRouter(newRouter)
		})
	}
}

//prepareRoutesNewRouter -
func (routeStructure RouteStructure) prepareRoutesNewRouter(router chi.Router) {
	router.Use(routeStructure.Uses...)

	for _, method := range routeStructure.Methods {
		router.Method(method.Method, method.Pattern, method.Handler)
	}

	for _, mountInfo := range routeStructure.Mounts {
		newRouter := chi.NewRouter()
		mountInfo.ServeFunc(newRouter, mountInfo.OtherParams...)
		router.Mount(mountInfo.Pattern, newRouter)
	}

	for pattern, routeInfo := range routeStructure.Routes {
		router.Route(pattern, func(newRouter chi.Router) {
			(*routeInfo).prepareRoutesNewRouter(newRouter)
		})
	}

	for _, groupInfo := range routeStructure.Groups {
		router.Group(func(newRouter chi.Router) {
			groupInfo.prepareRoutesNewRouter(newRouter)
		})
	}
}
