package bootstrap

import (
	"github.com/gin-gonic/gin"
	"gitlab.tiandy.com/lizewei08892/ginwebframework/global"
)

type Application struct {
	engine *gin.Engine
}

type StartOption func()

var startOpts = []StartOption{}

func HasStartOpt() bool {
	return len(startOpts) != 0
}

func StartOpt(opts ...StartOption) {
	startOpts = append(startOpts, opts...)
}

func doStartOpt() {
	if HasStartOpt() {
		for _, opt := range startOpts {
			opt()
		}
	}
}

func BootUp() {
	initAll()
	doStartOpt()
	global.App.Logger.Info("starting ------ ----- --- ")
	var ginEngine *gin.Engine
	if HasRouter() {
		ginEngine = SetupRouter()
	} else {
		ginEngine = gin.Default()
	}

	app := global.App.Config.App

	//服务启动
	_ = ginEngine.Run(":" + app.Port)
}

func initAll() {
	InitConfig()
	config := global.Config()
	if config == nil {
		return
	}
	InitLog()
	if !config.Database.IsEmpty() {
		InitDb()
	}
	if !config.Redis.IsEmpty() {
		InitRedis()
	}
}
