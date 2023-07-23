package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/healer1219/gin-web-framework/global"
)

type StartOption func()

type BootOption func() *global.Application

type Application struct {
	engine    *gin.Engine
	bootOpts  []BootOption
	startOpts []StartOption
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
	return &Application{
		engine:    gin.Default(),
		bootOpts:  opts,
		startOpts: make([]StartOption, 0),
		globalApp: global.App,
	}
}

func (app *Application) Router(opts ...RouteOption) *Application {
	for _, opt := range opts {
		opt(app.engine)
	}
	return app
}

func (app *Application) BootUp() *Application {
	for _, bootOpt := range app.bootOpts {
		bootOpt()
	}
	for _, startOpt := range app.startOpts {
		startOpt()
	}
	global.App.Logger.Info("starting ------ ----- --- ")
	_ = app.engine.Run(":" + app.globalApp.Config.App.Port)
	return app
}
