package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/healer1219/gin-web-framework/global"
)

type StartFunc func()

type BootOption func() *global.Application

type Application struct {
	engine    *gin.Engine
	bootOpts  []BootOption
	startOpts []StartFunc
	globalApp *global.Application
}

var baseBootOption = []BootOption{
	InitConfig,
	InitLog,
	InitDb,
}

func Default() *Application {
	return NewApplicationWithOpts(baseBootOption...)
}

func NewApplicationWithOpts(opts ...BootOption) *Application {
	return NewApplication(
		gin.Default(),
		opts,
		make([]StartFunc, 0),
		global.App,
	)
}

func NewApplication(engine *gin.Engine, bootOpts []BootOption, startOpts []StartFunc, globalApp *global.Application) *Application {
	return &Application{
		engine:    engine,
		bootOpts:  bootOpts,
		startOpts: startOpts,
		globalApp: globalApp,
	}
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

func (app *Application) BootUp() {
	for _, bootOpt := range app.bootOpts {
		bootOpt()
	}
	for _, startOpt := range app.startOpts {
		startOpt()
	}
	global.App.Logger.Info("starting ------ ----- --- ")
	_ = app.engine.Run(":" + app.globalApp.Config.App.Port)
}
