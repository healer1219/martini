package config

type Log struct {
	Level      string `mapstructure:"level" json:"level,omitempty" yaml:"level"`
	RootDir    string `mapstructure:"root_dir" json:"rootDir,omitempty" yaml:"rootDir"`
	FileName   string `mapstructure:"fileName" json:"fileName,omitempty" yaml:"fileName"`
	Format     string `mapstructure:"format" json:"format,omitempty" yaml:"format"`
	ShowLine   bool   `mapstructure:"showLine" json:"showLine,omitempty" yaml:"showLine"`
	MaxBackups int    `mapstructure:"maxBackups" json:"maxBackups,omitempty" yaml:"maxBackups"`
	MaxSize    int    `mapstructure:"maxSize" json:"maxSize,omitempty" yaml:"maxSize"`
	MaxAge     int    `mapstructure:"maxAge" json:"maxAge,omitempty" yaml:"maxAge"`
	Compress   bool   `mapstructure:"compress" json:"compress,omitempty" yaml:"compress"`
}
