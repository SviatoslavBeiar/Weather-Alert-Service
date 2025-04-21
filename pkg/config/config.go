package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBUser, DBPass, DBHost, DBPort, DBName string
	SMTPHost, SMTPPort, SMTPUser, SMTPPass string
}

func NewConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not found, using environment variables")
	}
	return Config{
		DBUser:   os.Getenv("DB_USER"),
		DBPass:   os.Getenv("DB_PASS"),
		DBHost:   os.Getenv("DB_HOST"),
		DBPort:   os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: os.Getenv("SMTP_PORT"),
		SMTPUser: os.Getenv("SMTP_USER"),
		SMTPPass: os.Getenv("SMTP_PASS"),
	}
}
