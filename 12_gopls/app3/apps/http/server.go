package main

import (
	"compress/flate"
	"context"
	"github.com/equilibristofgo/sandbox/04_internal/app3/apps/http/controller"
	_ "github.com/equilibristofgo/sandbox/04_internal/app3/docs"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/common/customRouter"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/common/properties"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/core/ports"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/core/service"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/infrastructure/persistence/mongoDB"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/infrastructure/persistence/postgreSQL"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"log"
	"net/http"
	_ "net/http/pprof"

	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	ApplicationProperties string = "./configs/application_properties.yaml"
	PostgresBBDDConfig    string = "./configs/bbdd_configuration_postgres.yaml"
	MongoBBDDConfig       string = "./configs/bbdd_configuration_mongo.yaml"
)

const (
	Postgres string = "postgres"
	Mongo    string = "mongo"
)

var (
	binder = Binder{}
)

// Binder is a set of services.
type Binder struct {
	router        *chi.Mux
	appProperties *properties.ApplicationProperties
	taskService   ports.TaskService
}

// @title           Swagger API micro-template
// @version         1.0
// @description     Micro-template is a microservice that has examples of database connection.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	appProperties := new(properties.ApplicationProperties)
	appProperties.GetConfiguration(ApplicationProperties)
	binder = NewBinder(appProperties)
	StartApp(appProperties)
}

//NewBinder - new initialization
func NewBinder(appProperties *properties.ApplicationProperties) Binder {
	var taskRepository ports.TasksRepository
	switch appProperties.Repository.Db {
	case Postgres:
		db := gorm.DB{} //TODO: create real db connection
		taskRepository = postgreSQL.NewTaskPostgres(&db)
	case Mongo:
		db := mongo.Client{} //TODO: create real db connection
		taskRepository = mongoDB.NewTaskMongoDB(&db)
	default:
		//panic("incorrect database")
	}

	taskService := service.NewTaskService(taskRepository)

	router := chi.NewRouter()
	return Binder{router,
		appProperties,
		taskService,
	}
}

//StartApp starts the application using the globally configured binder, and configuration
func StartApp(appProperties *properties.ApplicationProperties) {
	headerLog := "Main - StartApp - "
	binder.router = CORS(binder.router)
	binder.router = Routes(binder)
	log.Println(headerLog + appProperties.HttpHandler.Url + ":" + appProperties.HttpHandler.Port)
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Println(method + " " + route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(binder.router, walkFunc); err != nil {
		log.Println("Logging err: " + err.Error()) // panic if there is an error
	}
	binder.router.Get("/swagger/*", httpSwagger.WrapHandler)                                                                              // Expose Swagger documentation
	binder.router.Mount("/debug", middleware.Profiler())                                                                                  //debug options
	log.Println(headerLog + http.ListenAndServe(appProperties.HttpHandler.Url+":"+appProperties.HttpHandler.Port, binder.router).Error()) // Note, the port is usually gotten from the environment.
	panic(1)
}

//Routes configures the server routers
func Routes(binder Binder) *chi.Mux {
	rs := customRouter.RouteStructure{}

	compressor := middleware.NewCompressor(flate.DefaultCompression)
	rs.Uses = []func(http.Handler) http.Handler{render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,    // Log API request calls
		compressor.Handler,   // Compress results, mostly gzipping assets and json
		middleware.Recoverer, // Recover from panics without crashing server
	}
	taskRouter := chi.NewRouter()
	controller.NewTaskHttp(taskRouter, binder.taskService)

	rs.Routes = map[string]*customRouter.RouteStructure{
		"/v1": &customRouter.RouteStructure{
			Routes: map[string]*customRouter.RouteStructure{
				"/template": &customRouter.RouteStructure{
					Uses: []func(http.Handler) http.Handler{apiVersionCtx("v1")},
					Groups: []customRouter.RouteStructure{
						//Public

						//Private
						customRouter.RouteStructure{
							Mounts: []customRouter.RoutePatternInfo{
								customRouter.RoutePatternInfo{
									Pattern:     "/tasks",
									OtherParams: []interface{}{binder.taskService},
									ServeFunc: func(router *chi.Mux, otherParams ...interface{}) {
										controller.NewTaskHttp(router, otherParams[0].(ports.TaskService))
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return RoutesCustomizable(binder.router, rs)
}

func RoutesCustomizable(router *chi.Mux, routeStructure customRouter.RouteStructure) *chi.Mux {
	routeStructure.PrepareRoutes(router)
	return router
}

func CORS(router *chi.Mux) *chi.Mux {
	corsOptions := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Impersonation"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		//      Debug: true,
	})
	router.Use(corsOptions.Handler)
	return router
}

func apiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "api.version", version))
			next.ServeHTTP(w, r)
		})
	}
}
