package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DbFileName      string `yaml:"dbFileName"`
	TgBotToken      string `yaml:"tgBotToken"`
	UserTokenLength int    `yaml:"userTokenLength"`
	LogUserMessages bool   `yaml:"logUserMessages"`
	RestApiPort     int    `yaml:"restApiPort"`
	TlsCertFile     string `yaml:"tlsCertFile"`
	TlsKeyFile      string `yaml:"tlsKeyFile"`
	LogFileName     string `yaml:"logFileName"`
}

func Load(filename string) Config {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		log.Panic("Error opening config file:", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Panic("Error parsing config file:", err)
	}

	setDefaults(&config)
	validate(&config)

	return config
}

func setDefaults(config *Config) {
	if config.DbFileName == "" {
		config.DbFileName = "messageredir.db"
	}
	if config.UserTokenLength == 0 {
		config.UserTokenLength = 42
	}
	if config.RestApiPort == 0 {
		config.RestApiPort = 8083
	}
	if config.LogFileName == "" {
		config.LogFileName = "messageredir.log"
	}
}

func validate(config *Config) {
	if config.TgBotToken == "" {
		panic("TgBotToken is required")
	}
}
