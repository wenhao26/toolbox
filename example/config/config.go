package config

// viper把yaml的数据给对应的结构体

type ExampleConfig struct {
	Name     string `mapstructure:"name"`
	LogsPath string `mapstructure:"logs_path"`
	RedisCfg RedisConfig
	MySQLCfg MySQLConfig
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Db       string `mapstructure:"db"`
}
