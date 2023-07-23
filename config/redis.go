package config

import "reflect"

type Redis struct {
	Ip       string `json:"ip,omitempty" mapstructure:"ip" yaml:"ip"`
	Port     string `json:"port,omitempty" mapstructure:"port" yaml:"port"`
	Password string `json:"password,omitempty" mapstructure:"password" yaml:"password"`
	DbName   int    `json:"dbName" mapstructure:"db-name" yaml:"db-name"`
	PoolSize int    `json:"pool-size,omitempty" mapstructure:"pool-size" yaml:"pool-size"`
}

func (redis Redis) IsEmpty() bool {
	return reflect.DeepEqual(redis, Redis{})
}
