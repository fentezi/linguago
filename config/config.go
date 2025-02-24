package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Env      string   `yaml:"env"`
		Postgres Postgres `yaml:"postgres"`
		ApiKey   string   `yaml:"api_key" env:"API_KEY"`
	}

	Postgres struct {
		Host     string `yaml:"host" env:"POSTGRES_HOST"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
		Database string `yaml:"database"`
	}
)

func MustConfig() *Config {
	path := getPathFile()
	if path == " " {
		panic("config file path not provided")
	}

	if ok := fileExists(path); !ok {
		panic("config file does not exist")
	}

	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("failed to read config file: " + err.Error())
	}

	err = godotenv.Load()
	if err != nil {
		panic("failed to load environment variables: " + err.Error())
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic("failed to read environment variables: " + err.Error())
	}

	return &cfg

}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func getPathFile() string {
	var path string

	flag.StringVar(&path, "config", "path to the config file", "")

	flag.Parse()

	return path

}
