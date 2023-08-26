package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/healer1219/martini/global"
)

type StartFunc func()

type BootOption func() *global.Application

type Application struct {
	engine      *gin.Engine
	bootOpts    []BootOption
	startOpts   []StartFunc
	middleWares []gin.HandlerFunc
	globalApp   *global.Application
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

func (app *Application) Router(opts ...RouteOption) *Application {
	for _, opt := range opts {
		opt(app.engine)
	}
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

func (app *Application) BootUp() {
	for _, bootOpt := range app.bootOpts {
		bootOpt()
	}
	for _, middleWare := range app.middleWares {
		app.engine.Use(middleWare)
	}
	for _, startOpt := range app.startOpts {
		startOpt()
	}
	global.App.Logger.Info("starting ------ ----- --- ")
	_ = app.engine.Run(":" + app.globalApp.Config.App.Port)
}
