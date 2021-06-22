package configuration

import (
	"fmt"

	"github.com/joho/godotenv"
)

type ConfigService struct{}

func (ConfigService) Setup() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

func (ConfigService) Shutdown() {}
