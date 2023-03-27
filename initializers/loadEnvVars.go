package initializers

import (
	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	err := godotenv.Load()
	return err
}
