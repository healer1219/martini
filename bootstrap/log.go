package bootstrap

import (
	"github.com/healer1219/gin-web-framework/config"
	"github.com/healer1219/gin-web-framework/global"
	"github.com/healer1219/gin-web-framework/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

const (
	DebugLevel  = "debug"
	InfoLevel   = "info"
	WarnLevel   = "warn"
	ErrorLevel  = "error"
	DPanicLevel = "dpanic"
	PanicLevel  = "panic"
	FatalLevel  = "fatal"
)

var (
	level        zapcore.Level
	options      []zap.Option
	globalConfig *config.Config
)

func InitLog() *global.Application {
	global.App.RequireConfig(" init logger! ")
	globalConfig = global.App.Config
	createRootDir()
	setLogLevel()
	core := getZapCore()
	global.App.Logger = zap.New(core, options...)
	return global.App
}

func getZapCore() zapcore.Core {
	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(globalConfig.App.Env + "-" + l.String())
	}
	if strings.ToLower(globalConfig.Log.Format) == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	if globalConfig.Log.ShowLine {
		options = append(options, zap.AddCaller())
	}

	return zapcore.NewCore(encoder, logWriter(), level)

}

func logWriter() zapcore.WriteSyncer {
	writer := GetLogWriter(globalConfig.Log.FileName)
	return zapcore.AddSync(writer)
}

func GetLogWriter(logFileName string) *lumberjack.Logger {
	var log = globalConfig.Log
	writer := &lumberjack.Logger{
		// 文件名
		Filename: log.RootDir + string(os.PathSeparator) + logFileName,
		// 日志大小限制
		MaxSize: log.MaxSize,
		// 历史日志保留天数
		MaxAge: log.MaxAge,
		// 最大保留数量
		MaxBackups: log.MaxBackups,
		// 是否本地时区
		LocalTime: true,
		// 历史日志压缩标识
		Compress: false,
	}
	return writer
}

/** 创建日志文件夹 **/
func createRootDir() {
	rootDir := globalConfig.Log.RootDir
	if exists, _ := utils.PathExists(rootDir); !exists {
		_ = os.MkdirAll(rootDir, os.ModePerm)
	}
}

/** 日志级别转换 **/
func setLogLevel() {
	var lowerLevel string
	if globalConfig.Log.Level != "" {
		lowerLevel = strings.ToLower(globalConfig.Log.Level)
	}
	switch lowerLevel {
	case DebugLevel:
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(level))
	case InfoLevel:
		level = zap.InfoLevel
	case WarnLevel:
		level = zap.WarnLevel
	case ErrorLevel:
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(level))
	case DPanicLevel:
		level = zap.DPanicLevel
	case PanicLevel:
		level = zap.PanicLevel
	case FatalLevel:
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

}
