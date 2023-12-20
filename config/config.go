package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/healer1219/martini/global"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	App          App                    `mapstructure:"app" json:"app" yaml:"app"`
	Log          Log                    `mapstructure:"mlog" json:"mlog" yaml:"mlog"`
	Database     Database               `mapstructure:"database" json:"database" yaml:"database"`
	DatabaseMap  map[string]Database    `mapstructure:"dbs" json:"dbs" yaml:"dbs"`
	Redis        Redis                  `mapstructure:"redis" json:"redis" yaml:"redis"`
	CustomConfig map[string]interface{} `mapstructure:"custom" json:"custom" yaml:"custom"`
	Cloud        Registry               `mapstructure:"cloud" json:"cloud" yaml:"cloud"`
}

func InitConfig() *global.Application {
	configFile := "config.yaml"
	if envConfigFile := os.Getenv("CONFIG_FILE"); envConfigFile != "" {
		configFile = envConfigFile
	}
	InitConfigByName(configFile)
	return global.App
}

func InitConfigByName(fileName string) *viper.Viper {
	v := viper.New()
	v.SetConfigFile(fileName)
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

func GetConfig(configFileName string) (*viper.Viper, Config) {
	v := viper.New()
	v.SetConfigFile(configFileName)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("read config file failed: %s \n", err))
	}
	conf := Config{}
	v.Unmarshal(&conf)
	return v, conf
}
