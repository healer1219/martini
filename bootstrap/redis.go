package bootstrap

import (
	"github.com/go-redis/redis"
	"github.com/healer1219/martini/config"
	"github.com/healer1219/martini/global"
	"go.uber.org/zap"
)

func InitRedis() *global.Application {
	global.App.RequireConfigAndLog(" init Redis! ")
	client, err := getRedisClient(&global.App.Config.Redis)
	if err != nil {
		panic("init redis failed!")
	}
	global.App.RedisClient = client
	return global.App
}

func getRedisClient(clientConf *config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     clientConf.Ip + ":" + clientConf.Port,
		DB:       clientConf.DbName,
		Password: clientConf.Password,
		PoolSize: clientConf.PoolSize,
	})
	_, err := client.Ping().Result()
	if err != nil {
		global.Logger().Fatal("redis connect failed!", zap.Any("err", err))
	}
	return client, nil
}
