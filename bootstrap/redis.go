package bootstrap

import (
	"github.com/go-redis/redis"
	"github.com/healer1219/gin-web-framework/config"
	"github.com/healer1219/gin-web-framework/global"
	"go.uber.org/zap"
)

func InitRedis() (*redis.Client, error) {
	return getRedisClient(&global.App.Config.Redis)
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
