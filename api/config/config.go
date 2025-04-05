package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	Port     string
	ApiQuota string
	Domain   string
	DbUrl    string
	DbPass   string
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	Port = os.Getenv("PORT")
	ApiQuota = os.Getenv("API_QUOTA")
	Domain = os.Getenv("DOMAIN")
	DbUrl = os.Getenv("DB_ADD")
	DbPass = os.Getenv("DB_PASS")
}
