// Package config main config
package config

import (
	"fmt"

	"github.com/caarlos0/env/v7"
)

// MainConfig with init data
type MainConfig struct {
	PostgresPort     string `env:"POSTGRES_PORT,notEmpty" envDefault:"5432"`
	PostgresHost     string `env:"POSTGRES_HOST,notEmpty" envDefault:"localhost"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,notEmpty" envDefault:"postgres"`
	PostgresUser     string `env:"POSTGRES_USER,notEmpty" envDefault:"postgres"`
	PostgresDB       string `env:"POSTGRES_DB,notEmpty" envDefault:"postgres"`
	JwtKey           string `env:"JWT_KEY,notEmpty" envDefault:"874967EC3EA3490F8F2EF6478B72A756"`
	Port             string `env:"PORT,notEmpty" envDefault:"10000"`
	Host             string `env:"HOST,notEmpty" envDefault:"localhost"`
}

// NewMainConfig parsing config from environment
func NewMainConfig() (*MainConfig, error) {
	mainConfig := &MainConfig{}

	err := env.Parse(mainConfig)
	if err != nil {
		return nil, fmt.Errorf("config - NewMainConfig - Parse:%w", err)
	}

	return mainConfig, nil
}
