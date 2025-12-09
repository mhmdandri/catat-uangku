package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret   string
	JWTIssuer   string
	JWTAudience string

	GoogleClientID      string
	GoogleClientSecret  string
	GoogleRedirectURL   string
	FrontendRedirectURL string
}

var Cfg *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env tidak ada, silahkan cek kembali")
	}
	Cfg = &Config{
		AppPort:    os.Getenv("APP_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),

		JWTSecret:   os.Getenv("JWT_SECRET"),
		JWTIssuer:   os.Getenv("JWT_ISSUER"),
		JWTAudience: os.Getenv("JWT_AUDIENCE"),

		GoogleClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:  os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:   os.Getenv("GOOGLE_REDIRECT_URL"),
		FrontendRedirectURL: os.Getenv("FRONTEND_REDIRECT_URL"),
	}
}
