package config

import "github.com/spf13/viper"

type Config struct {
	Server struct {
		Address string
	}
	Redis struct {
		URL      string
		Password string
		DB       int
	}
}

func Load() *Config {
	// Defaults
	viper.SetDefault("SERVER_ADDR", ":8080")
	viper.SetDefault("REDIS_ADDR", "redis:6379")

	var cfg Config
	cfg.Server.Address = viper.GetString("SERVER_ADDR")

	cfg.Redis.URL = viper.GetString("REDIS_ADDR")
	cfg.Redis.Password = viper.GetString("REDIS_PASSWORD")
	cfg.Redis.DB = viper.GetInt("REDIS_DB")

	return &cfg
}
