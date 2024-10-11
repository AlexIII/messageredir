package config

import (
	"log"
	"os"

	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v2"
)

const (
	EnvConfigPrefix = "MREDIR_"
)

type Config struct {
	DbFileName      string `yaml:"dbFileName" env:"DB_FILE_NAME"`
	TgBotToken      string `yaml:"tgBotToken" env:"TG_BOT_TOKEN"`
	UserTokenLength int    `yaml:"userTokenLength" env:"USER_TOKEN_LENGTH"`
	LogUserMessages bool   `yaml:"logUserMessages" env:"LOG_USER_MESSAGES"`
	RestApiPort     int    `yaml:"restApiPort" env:"REST_API_PORT"`
	TlsCertFile     string `yaml:"tlsCertFile" env:"TLS_CERT_FILE"`
	TlsKeyFile      string `yaml:"tlsKeyFile" env:"TLS_KEY_FILE"`
	LogFileName     string `yaml:"logFileName" env:"LOG_FILE_NAME"`
}

func loadFromYaml(filename string) Config {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Cannot opening config file:", err, "Skipping.")
		return config
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Panic("Error parsing config file:", err)
	}

	return config
}

func enrichFromEnv(config *Config) {
	env.ParseWithOptions(config, env.Options{Prefix: EnvConfigPrefix}) // Ignore errors
}

func Load(filename string) Config {
	config := loadFromYaml(filename)
	enrichFromEnv(&config)
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
}

func validate(config *Config) {
	if config.TgBotToken == "" {
		panic("TgBotToken is required")
	}
}
