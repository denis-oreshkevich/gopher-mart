package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
)

const (
	serverAddressEnvName = "RUN_ADDRESS"

	dataBaseURIEnvName          = "DATABASE_URI"
	accrualSystemAddressEnvName = "ACCRUAL_SYSTEM_ADDRESS"

	defaultHost = "localhost"

	defaultPort = "8081"

	defaultAccrualSystemAddress = "http://localhost:8080"
)

var config Config

func Get() Config {
	return config
}

func Parse() (*Config, error) {
	flag.StringVar(&config.serverAddress, "a", fmt.Sprintf("%s:%s", defaultHost, defaultPort),
		"HTTP server address")
	flag.StringVar(&config.dataBaseURI, "d",
		"host=localhost port=5433 user=postgres password=postgres dbname=mart sslmode=disable",
		"DataBase URI")
	flag.StringVar(&config.accrualSystemAddress, "r", defaultAccrualSystemAddress,
		"Accrual System Address")

	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("env.Parse: %w", err)
	}

	logger.Log.Info(fmt.Sprintf("initializing config %+v", config))

	return &config, nil
}
