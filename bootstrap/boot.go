package bootstrap

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/healer1219/martini/cloud"
	"github.com/healer1219/martini/global"
	"log"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type StartFunc func()
type ShutDownFunc func()

type BootOption func() *global.Application

type Application struct {
	engine          *gin.Engine
	bootOpts        []BootOption
	startOpts       []StartFunc
	shutDownOpts    []ShutDownFunc
	middleWares     []gin.HandlerFunc
	globalApp       *global.Application
	serviceInstance cloud.ServiceInstance
	registry        cloud.ServiceRegistry
}

var baseBootOption = []BootOption{
	InitConfig,
	InitLog,
}

func Default() *Application {
	app := NewApplicationWithOpts()
	return app
}

func NewApplicationWithOpts(opts ...BootOption) *Application {
	return NewApplication(
		newGin(),
		opts,
		make([]StartFunc, 0),
		global.App,
	)
}

func newGin() *gin.Engine {
	for _, bootOpt := range baseBootOption {
		bootOpt()
	}
	engine := gin.New()
	engine.Use(
		LoggerMiddleWare(global.Logger()),
		GinRecovery(global.Logger(), true),
	)
	return engine
}

func NewApplication(engine *gin.Engine, bootOpts []BootOption, startOpts []StartFunc, globalApp *global.Application) *Application {
	return &Application{
		engine:    engine,
		bootOpts:  bootOpts,
		startOpts: startOpts,
		globalApp: globalApp,
	}
}

func (app *Application) BootOpt(bootOpts ...BootOption) *Application {
	if app.bootOpts == nil {
		app.bootOpts = bootOpts
	} else {
		app.bootOpts = append(app.bootOpts, bootOpts...)
	}
	return app
}

func (app *Application) StartFunc(startOpts ...StartFunc) *Application {
	if app.startOpts == nil {
		app.startOpts = startOpts
	} else {
		app.startOpts = append(app.startOpts, startOpts...)
	}
	return app
}

func (app *Application) ShutDownFunc(shutDownOpts ...ShutDownFunc) *Application {
	if app.shutDownOpts == nil {
		app.shutDownOpts = shutDownOpts
	} else {
		app.shutDownOpts = append(app.shutDownOpts, shutDownOpts...)
	}
	return app
}

func (app *Application) Router(opts ...RouteOption) *Application {
	Regist(opts...)
	return app
}

func (app *Application) Use(middleware ...gin.HandlerFunc) *Application {
	if app.middleWares == nil {
		app.middleWares = middleware
	} else {
		app.middleWares = append(app.middleWares, middleware...)
	}
	return app
}

func (app *Application) Discovery(serviceInstance cloud.ServiceInstance, registry cloud.ServiceRegistry) *Application {
	app.serviceInstance = serviceInstance
	app.registry = registry
	app.StartFunc(func() {
		registry.Register(serviceInstance)
	})
	app.ShutDownFunc(func() {
		registry.Deregister()
	})
	return app
}

func (app *Application) DefaultDiscovery() *Application {
	instance, err := cloud.NewDefaultServiceInstance()
	if err != nil {
		log.Fatal(err)
	}

	registry := global.Config().Cloud
	if registry.IsEmpty() {
		log.Fatal("config file [cloud] is illegal")
	}
	serviceRegistry, err := cloud.NewDefaultConsulServiceRegistry(registry.Ip, registry.Port, registry.Token)
	if err != nil {
		log.Fatal(err)
	}

	app.Router(func(engine *gin.Engine) {
		engine.GET("/actuator/health", cloud.DefaultHealthCheck)
	})

	return app.Discovery(instance, serviceRegistry)
}

func (app *Application) BootUp() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	for _, bootOpt := range app.bootOpts {
		bootOpt()
	}
	for _, middleWare := range app.middleWares {
		app.engine.Use(middleWare)
	}
	for _, opt := range routeOpts {
		opt(app.engine)
	}
	for _, startOpt := range app.startOpts {
		startOpt()
	}
	global.App.Logger.Info("starting ------ ----- --- ")
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(app.globalApp.Config.App.Port),
		Handler: app.engine,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("application run failed!", err)
		}
	}()

	<-ctx.Done()
	for _, shutDownOpt := range app.shutDownOpts {
		shutDownOpt()
	}
	stop()
	log.Println("application is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("application forced to shutdown: ", err)
	}

	log.Println("application exiting")
}
