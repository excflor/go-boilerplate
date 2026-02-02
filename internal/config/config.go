package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppPort string `env:"APP_PORT" env-default:"8080"`

	Database struct {
		Host     string `env:"DB_HOST" env-required:"true"`
		Port     string `env:"DB_PORT" env-required:"true"`
		User     string `env:"DB_USER" env-required:"true"`
		Password string `env:"DB_PASSWORD" env-required:"true"`
		Name     string `env:"DB_NAME" env-required:"true"`

		MaxOpenConns    int `env:"DB_MAX_OPEN_CONNS" env-default:"25"`
		MaxIdleConns    int `env:"DB_MAX_IDLE_CONNS" env-default:"10"`
		ConnMaxLifetime int `env:"DB_CONN_MAX_LIFETIME" env-default:"15"` // in minutes
	}

	MaxRequestPerSecond float64 `env:"MAX_REQUEST_PER_SECOND" env-default:"20"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(".env", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
