package configuration

import "github.com/joho/godotenv"

type ConfigService struct{}

func (ConfigService) Setup() error {
	return godotenv.Load()
}

func (ConfigService) Shutdown() {}
