package global

import (
	"github.com/go-redis/redis"
	"github.com/healer1219/gin-web-framework/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application struct {
	ConfigViper *viper.Viper
	Config      *config.Config
	Logger      *zap.Logger
	DB          *gorm.DB
	RedisClient *redis.Client
	dbMap       map[string]*gorm.DB
}

var App = new(Application)

func Config() *config.Config {
	return App.Config
}

func GetConfigByName(key string) interface{} {
	return App.Config.CustomConfig[key]
}

func Logger() *zap.Logger {
	return App.Logger
}

func DB() *gorm.DB {
	return App.DB
}

func (app *Application) RequireConfig(msg string) {
	if app.Config == nil {
		panic("config is nil, msg: " + msg)
	}
}

func (app *Application) RequireLog(msg string) {
	if app.Logger == nil {
		panic("logger is nil, msg: " + msg)
	}
}

func (app *Application) RequireConfigAndLog(msg string) {
	if app.Config == nil {
		panic("config and logger are both nil, msg: " + msg)
	}
}

func (app *Application) DbByName(customName string) *gorm.DB {
	if app.dbMap == nil || len(app.dbMap) == 0 {
		app.Logger.Warn("dataBase connections is empty")
		return app.DB
	}
	return app.dbMap[customName]
}

func (app *Application) AddDb(customName string, db *gorm.DB) *gorm.DB {
	if app.dbMap == nil {
		app.dbMap = make(map[string]*gorm.DB)
	}
	app.dbMap[customName] = db
	return app.dbMap[customName]
}

func (app *Application) Redis() *redis.Client {
	return app.RedisClient
}
