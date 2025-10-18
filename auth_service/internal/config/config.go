package config

import (
	"os"
)

type Config struct {
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string

	SERV_HOST string
	SERV_PORT string
}

func LoadConfig() (*Config, error) {
	return &Config{
		DB_HOST:     os.Getenv("AUTH_DB_HOST"),
		DB_PORT:     os.Getenv("AUTH_DB_PORT"),
		DB_USER:     os.Getenv("AUTH_DB_USER"),
		DB_PASSWORD: os.Getenv("AUTH_DB_PASSWORD"),
		DB_NAME:     os.Getenv("AUTH_DB_NAME"),
		SERV_HOST:   os.Getenv("AUTH_SERV_HOST"),
		SERV_PORT:   os.Getenv("AUTH_SERV_PORT"),
	}, nil
}
