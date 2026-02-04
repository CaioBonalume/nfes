package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	CERT_PATH string `env:"CERT_PATH"`
	CERT_PASS string `env:"KEY_PATH"`
	CERT_PEM  string `env:"CERT_PEM"`
	CERT_KEY  string `env:"CERT_KEY"`
}

var Env env

func LoadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	Env = env{
		CERT_PATH: os.Getenv("CERT_PATH"),
		CERT_PASS: os.Getenv("KEY_PATH"),
		CERT_PEM:  os.Getenv("CERT_PEM"),
		CERT_KEY:  os.Getenv("CERT_KEY"),
	}
}
