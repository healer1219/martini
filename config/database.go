package config

import "reflect"

type Database struct {
	DatabaseType        string `mapstructure:"database_type"  json:"database_type,omitempty" yaml:"database_type"`
	Ip                  string `mapstructure:"ip" json:"ip,omitempty" yaml:"ip"`
	Port                int    `mapstructure:"port" json:"port,omitempty" yaml:"port"`
	DatabaseName        string `mapstructure:"database_name" json:"database_name,omitempty" yaml:"database_name"`
	UserName            string `mapstructure:"username" json:"username,omitempty" yaml:"username"`
	Password            string `mapstructure:"password" json:"password,omitempty" yaml:"password"`
	Charset             string `mapstructure:"charset" json:"charset,omitempty" yaml:"charset"`
	MaxIdleConns        int    `mapstructure:"max_idle_conns" json:"max_idle_conns,omitempty" yaml:"max_idle_conns"`
	MaxOpenConns        int    `mapstructure:"max_open_conns" json:"max_open_conns,omitempty" yaml:"max_open_conns"`
	LogMode             string `mapstructure:"log_mode" json:"log_mode,omitempty" yaml:"log_mode"`
	EnableFileLogWriter bool   `mapstructure:"enable_file_log_writer" json:"enable_file_log_writer,omitempty" yaml:"enable_file_log_writer"`
	LogFileName         string `mapstructure:"log_filename" json:"log_filename,omitempty" yaml:"log_filename"`
}

func (database Database) IsEmpty() bool {
	return reflect.DeepEqual(database, Database{})
}
