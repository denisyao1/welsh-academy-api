package common

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Configuration struct {
	DB_HOST     string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
	DB_PORT     int
	JWT_SECRET  string
}

func LoadConfig() Configuration {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load env variables")
	}

	var config Configuration
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))

	if err != nil {
		log.Fatal("Failed to parsed database port")
	}

	config.DB_HOST = os.Getenv("DB_HOST")
	config.DB_NAME = os.Getenv("DB_NAME")
	config.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	config.DB_USER = os.Getenv("DB_USER")
	config.DB_PORT = port
	config.JWT_SECRET = os.Getenv("JWT_SECRET")
	return config
}
