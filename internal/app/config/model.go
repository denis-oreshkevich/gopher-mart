package config

type Config struct {
	serverAddress        string
	dataBaseURI          string
	accrualSystemAddress string
}

func (c *Config) ServerAddress() string {
	return c.serverAddress
}

func (c *Config) DataBaseURI() string {
	return c.dataBaseURI
}

func (c *Config) AccrualSystemAddress() string {
	return c.accrualSystemAddress
}
