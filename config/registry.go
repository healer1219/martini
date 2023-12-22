package config

import (
	"reflect"
)

type Registry struct {
	Ip              string `mapstructure:"ip" json:"ip,omitempty" yaml:"ip"`
	Port            int    `mapstructure:"port" json:"port,omitempty" yaml:"port"`
	Token           string `mapstructure:"token" json:"token,omitempty" yaml:"token"`
	HealthCheckPath string `mapstructure:"health_check_path" json:"health_check_path,omitempty" yaml:"health_check_path"`
}

func (registry Registry) IsEmpty() bool {
	return reflect.DeepEqual(registry, Registry{})
}
