package bootstrap

import (
	"github.com/healer1219/martini/config"
	"github.com/healer1219/martini/global"
	"github.com/healer1219/martini/mlog"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	SilentMode = "silent"
	ErrorMode  = "error"
	WarnMode   = "warn"
	InfoMode   = "info"
)

func GetDb(dbConfig config.Database) *gorm.DB {
	switch dbConfig.DatabaseType {
	case "mysql":
		return initMysql(dbConfig)
	default:
		return initMysql(dbConfig)
	}
}

func InitDb() *global.Application {
	global.App.RequireConfigAndLog("init db!")
	dbConfig := global.App.Config.Database
	global.App.DB = GetDb(dbConfig)
	dbConfMap := global.Config().DatabaseMap
	if dbConfMap != nil && len(dbConfMap) != 0 {
		for k, conf := range dbConfMap {
			db := GetDb(conf)
			global.App.AddDb(k, db)
		}
	}
	return global.App
}

func RealeaseDB() {
	if global.App.DB != nil {
		db, _ := global.App.DB.DB()
		_ = db.Close()
	}
}

func initMysql(dbConfig config.Database) *gorm.DB {
	if dbConfig.DatabaseName == "" {
		return nil
	}
	dsn := generateMysqlUrl(dbConfig)
	mysqlConf := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         191,
		DisableDatetimePrecision:  true,
		DontSupportRenameColumn:   true,
		DontSupportRenameIndex:    true,
		SkipInitializeWithVersion: false,
	}
	db, err := gorm.Open(mysql.New(mysqlConf), &gorm.Config{
		//禁用自动创建外键
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   getGormLogger(dbConfig),
	})
	if err != nil {
		global.Logger().Error("mysql connect failed, err:", zap.Any("err", err))
		return nil
	}
	gormDb, _ := db.DB()
	gormDb.SetMaxIdleConns(dbConfig.MaxIdleConns)
	gormDb.SetMaxOpenConns(dbConfig.MaxOpenConns)
	return db
}

func generateMysqlUrl(conf config.Database) string {
	return conf.UserName + ":" + conf.Password + "@tcp(" + conf.Ip + ":" + strconv.Itoa(conf.Port) + ")/" + conf.DatabaseName + "?charset=" + conf.Charset + "&parseTime=True&loc=Local"
}

func getGormLogWriter(dbConfig config.Database) logger.Writer {
	var writer io.Writer
	if dbConfig.EnableFileLogWriter {
		writer = mlog.GetLogWriter(dbConfig.LogFileName)
	} else {
		writer = os.Stdout
	}
	return log.New(writer, "\r\n", log.LstdFlags)
}

func getGormLogger(dbConfig config.Database) logger.Interface {
	var logMode logger.LogLevel

	switch strings.ToLower(dbConfig.LogMode) {
	case ErrorMode:
		logMode = logger.Error
	case SilentMode:
		logMode = logger.Silent
	case WarnMode:
		logMode = logger.Warn
	case InfoMode:
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	return logger.New(getGormLogWriter(dbConfig), logger.Config{
		SlowThreshold:             3000 * time.Millisecond,
		LogLevel:                  logMode,
		IgnoreRecordNotFoundError: false,
		Colorful:                  !dbConfig.EnableFileLogWriter,
	})
}
