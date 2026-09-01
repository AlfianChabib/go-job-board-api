package config

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	DatabaseUser     string
	DatabasePassword string
	DatabaseDB       string
	DatabaseHost     string
	DatabasePort     int
	DatabaseUrl      string
	Port             string
}

func LoadEnv() (*Env, error) {
	config := viper.New()
	config.SetConfigFile(".env")

	err := config.ReadInConfig()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	var env *Env = &Env{
		DatabaseUser:     config.GetString("DATABASE_USER"),
		DatabasePassword: config.GetString("DATABASE_PASSWORD"),
		DatabaseDB:       config.GetString("DATABASE_DB"),
		DatabaseHost:     config.GetString("DATABASE_HOST"),
		DatabasePort:     config.GetInt("DATABASE_PORT"),
		DatabaseUrl:      config.GetString("DATABASE_URL"),
		Port:             config.GetString("PORT"),
	}

	return env, nil
}
