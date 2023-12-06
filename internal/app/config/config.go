package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
)

const (
	defaultHost = "localhost"

	defaultPort = "8081"

	defaultAccrualSystemAddress = "http://localhost:8080"
)

var conf Config

func Get() Config {
	return conf
}

func Parse() (*Config, error) {
	flag.StringVar(&conf.ServerAddress, "a", fmt.Sprintf("%s:%s", defaultHost, defaultPort),
		"HTTP server address")
	flag.StringVar(&conf.DataBaseURI, "d",
		"host=localhost port=5433 user=postgres password=postgres dbname=mart sslmode=disable",
		"DataBase URI")
	flag.StringVar(&conf.AccrualSystemAddress, "r", defaultAccrualSystemAddress,
		"Accrual System Address")

	flag.Parse()
	err := env.Parse(&conf)
	if err != nil {
		return nil, fmt.Errorf("env.Parse: %w", err)
	}

	logger.Log.Info(fmt.Sprintf("initializing Config %+v", conf))

	return &conf, nil
}
