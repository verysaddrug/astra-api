package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	AdminToken string
}

func LoadConfig(envFile string) *Config {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("Не удалось загрузить файл окружения %s: %v", envFile, err)
	}
	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		AdminToken: os.Getenv("ADMIN_TOKEN"),
	}
}
