package bootstrap

import "github.com/gin-gonic/gin"

type Option func(*gin.Engine)

var routeOpts = []Option{}

func HasRouter() bool {
	return len(routeOpts) != 0
}

func Regist(opts ...Option) {
	routeOpts = append(routeOpts, opts...)
}

func SetupRouter() *gin.Engine {
	engine := gin.Default()
	for _, routeOpt := range routeOpts {
		routeOpt(engine)
	}

	return engine
}
