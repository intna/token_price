package utils

import (
	"log"
	"strconv" //data type conv

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading .env file: %v", err)
	} else {
		log.Println(".env file loaded successfully")
	}
}

func GetValue[T any](key string, defaultValue T) T {
	if !viper.IsSet(key) {
		log.Printf("Key '%s' not found. Using default value: %v", key, defaultValue)
		return defaultValue
	}

	value := viper.GetString(key)

	var result T
	switch any(result).(type) {
	case int:
		parsedValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Failed to convert key '%s' value to int. Using default: %v", key, defaultValue)
			return defaultValue
		}
		return any(parsedValue).(T)
	case bool:
		parsedValue, err := strconv.ParseBool(value)
		if err != nil {
			log.Printf("Failed to convert key '%s' value to bool. Using default: %v", key, defaultValue)
			return defaultValue
		}
		return any(parsedValue).(T)
	case float64:
		parsedValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Printf("Failed to convert key '%s' value to float64. Using default: %v", key, defaultValue)
			return defaultValue
		}
		return any(parsedValue).(T)
	case string:
		return any(value).(T)
	default:
		log.Printf("Unsupported type for key '%s'. Using default: %v", key, defaultValue)
		return defaultValue
	}
}
