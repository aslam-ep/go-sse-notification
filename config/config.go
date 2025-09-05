package config

import "github.com/spf13/viper"

type Config struct {
	Server struct {
		Address string
	}
	Redis struct {
		URL string
	}
}

func Load() *Config {
	var cfg Config
	cfg.Server.Address = viper.GetString("server.address")
	cfg.Redis.URL = viper.GetString("redis.url")
	return &cfg
}
