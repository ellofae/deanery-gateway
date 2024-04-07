package config

import (
	"os"
	"path/filepath"

	"github.com/ellofae/deanery-gateway/pkg/environment"
	"github.com/ellofae/deanery-gateway/pkg/logger"
	"github.com/spf13/viper"
)

var configuration_path string = "./config"

type Config struct {
	ServerSettings struct {
		BindAddr     string `yaml:"bindAddr"`
		ReadTimeout  string `yaml:"readTimeout"`
		WriteTimeout string `yaml:"writeTimeout"`
		IdleTimeout  string `yaml:"idleTimeout"`
	} `yaml:"ServerSettings"`

	Authentication struct {
		JwtSecretToken string `yaml:"jwtSecretToken"`
	} `yaml:"Authentication"`

	Session struct {
		SessionKey string `yaml:"sessionKey"`
	} `yaml:"Session"`
}

func ConfigureViper() *viper.Viper {
	logger := logger.GetLogger()

	configurationFile, err := environment.ParseEnvironmentVariable()
	if err != nil {
		os.Exit(1)
	}

	filepath := filepath.Join(configuration_path, configurationFile)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		logger.Printf("Configuration file '%s' doesn't exist.\n", filepath)
		os.Exit(1)
	}

	v := viper.New()
	v.AddConfigPath(configuration_path)
	v.SetConfigName(configurationFile)
	v.SetConfigType("yaml")

	err = v.ReadInConfig()
	if err != nil {
		logger.Printf("Unable to read the configuration file. Error: %v.\n", err.Error())
		os.Exit(1)
	}
	logger.Println("Config loaded successfully.")

	return v
}

func ParseConfig(v *viper.Viper) *Config {
	logger := logger.GetLogger()

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		logger.Println("Unable to parse the configuration file.")
	}
	logger.Println("Configuratin file parsed successfully.")

	return cfg
}
