package routes

import "github.com/gin-gonic/gin"

type RouteOption func(*gin.Engine)

var routeOpts = []RouteOption{}

func HasRouter() bool {
	return len(routeOpts) != 0
}

func Regist(opts ...RouteOption) {
	routeOpts = append(routeOpts, opts...)
}

func SetupRouter() *gin.Engine {
	engine := gin.Default()
	for _, routeOpt := range routeOpts {
		routeOpt(engine)
	}

	return engine
}
