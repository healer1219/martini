package config

type App struct {
	Name string `mapstructure:"name" json:"name"  yaml:"name"`
	Env  string `mapstructure:"env" json:"env" yaml:"env"`
	Port string `mapstructure:"port" json:"port" yaml:"port"`
}
