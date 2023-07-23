package bootstrap

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gitlab.tiandy.com/lizewei08892/ginwebframework/config"
	"gitlab.tiandy.com/lizewei08892/ginwebframework/global"
	"os"
)

func InitConfig() *viper.Viper {
	configFile := "config.yaml"
	if envConfigFile := os.Getenv("CONFIG_FILE"); envConfigFile != "" {
		configFile = envConfigFile
	}
	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("read config file failed: %s \n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		UnmarshalConfig(v)
	})
	UnmarshalConfig(v)
	return v
}

func UnmarshalConfig(v *viper.Viper) {
	unmarshalErr := v.Unmarshal(&global.App.Config)
	if unmarshalErr != nil {
		fmt.Println(unmarshalErr)
	}
}

func GetConfig(configFileName string) (*viper.Viper, config.Config) {
	v := viper.New()
	v.SetConfigFile(configFileName)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("read config file failed: %s \n", err))
	}
	conf := config.Config{}
	v.Unmarshal(&conf)
	return v, conf
}
