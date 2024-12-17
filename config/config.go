package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Env      string   `yml:"env"`
		Redis    Redis    `yml:"redis"`
		Postgres Postgres `yml:"postgres"`
	}

	Redis struct {
		Host     string `yml:"host"`
		Port     string `yml:"port"`
		Password string `yml:"password"`
		NumberDB int    `yml:"db"`
	}

	Postgres struct {
		Host     string `yml:"host"`
		Port     string `yml:"port"`
		Username string `yml:"username"`
		Password string `yml:"password"`
		Database string `yml:"database"`
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