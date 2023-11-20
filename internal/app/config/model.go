package config

type Config struct {
	ServerAddress        string `env:"RUN_ADDRESS"`
	DataBaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}
