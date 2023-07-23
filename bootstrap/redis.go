package bootstrap

import (
	"github.com/go-redis/redis"
	"gitlab.tiandy.com/lizewei08892/ginwebframework/config"
	"gitlab.tiandy.com/lizewei08892/ginwebframework/global"
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
