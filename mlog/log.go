package mlog

import (
	"github.com/gin-gonic/gin"
	"github.com/healer1219/martini/config"
	"github.com/healer1219/martini/global"
	"github.com/healer1219/martini/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
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
	log.Printf("initing log \n")
	globalConfig = global.App.Config
	createRootDir()
	setLogLevel()
	core := getZapCore()
	global.App.Logger = zap.New(core, options...)
	zap.ReplaceGlobals(global.App.Logger)
	log.Printf("init log complete \n")
	return global.App
}

func LoggerMiddleWare(logger *zap.Logger) gin.HandlerFunc {
	global.App.RequireLog("init logger middle ware failed, require logger! ")
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)
		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		logger.Info(path,
			zap.Int("status", param.StatusCode),
			zap.String("method", param.Method),
			zap.String("path", path),
			zap.String("query", raw),
			zap.String("ip", param.ClientIP),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", param.ErrorMessage),
			zap.Duration("cost", param.Latency),
		)
	}
}

func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	global.App.RequireLog("init recovery failed, require logger! ")
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
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
