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
		DB_HOST:     os.Getenv("BOOKING_DB_HOST"),
		DB_PORT:     os.Getenv("BOOKING_DB_PORT"),
		DB_USER:     os.Getenv("BOOKING_DB_USER"),
		DB_PASSWORD: os.Getenv("BOOKING_DB_PASSWORD"),
		DB_NAME:     os.Getenv("BOOKING_DB_NAME"),
		SERV_HOST:   os.Getenv("BOOKING_SERV_HOST"),
		SERV_PORT:   os.Getenv("BOOKING_SERV_PORT"),
	}, nil
}
