package config

import (
	"github.com/caarlos0/env/v8"
)

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	Host string `env:"DB_HOST,required"`
	Port int    `env:"DB_PORT,required"`
	User string `env:"DB_USER,required"`
	Pass string `env:"DB_PASSWORD,required"`
	Name string `env:"DB_NAME,required"`
}

func (cfg *Config) ParseFromEnv() error {
	err := env.Parse(cfg)
	return err
}
