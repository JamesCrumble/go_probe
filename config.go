package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func ApiAddrFromValues(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

type Config struct {
	API_HOST    string `env:"API_HOST" envDefault:"0.0.0.0"`
	API_PORT    int    `env:"API_PORT" envDefault:"4200"`
	PROF_ENABLE bool   `env:"PROF_ENABLE" envDefault:"False"`
	PROF_HOST   string `env:"PROF_HOST" envDefault:"0.0.0.0"`
	PROF_PORT   int    `env:"PROF_PORT" envDefault:"8585"`
}

func (config *Config) ApiAddr() string {
	return ApiAddrFromValues(config.API_HOST, config.API_PORT)
}

func (config *Config) ProfAddr() string {
	return ApiAddrFromValues(config.API_HOST, config.API_PORT)
}

func initConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var config Config
	err = env.Parse(&config)
	if err != nil {
		panic(err)
	}

	return &config
}

var config = initConfig()
