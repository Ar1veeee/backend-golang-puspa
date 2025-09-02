package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	wd, _ := os.Getwd()
	log.Printf("Current working directory: %s", wd)
	
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Println("Info: No .env file found, using system environment variables")
		return
	}
	
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
		log.Println("Using system environment variables")
	} else {
		log.Println("✅ .env file loaded successfully")
	}
}

func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
