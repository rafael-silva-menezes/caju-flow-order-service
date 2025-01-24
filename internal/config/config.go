package config

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Conf struct {
	DBDriver             string `mapstructure:"DB_DRIVER"`
	DBHost               string `mapstructure:"DB_HOST"`
	DBPort               string `mapstructure:"DB_PORT"`
	DBUser               string `mapstructure:"DB_USER"`
	DBPassword           string `mapstructure:"DB_PASSWORD"`
	DBName               string `mapstructure:"DB_NAME"`
	BrokerPort           string `mapstructure:"BROKER_PORT"`
	BrokerManagementPort string `mapstructure:"BROKER_MANAGEMENT_PORT"`
	BrokerUser           string `mapstructure:"BROKER_USER"`
	BrokerPassword       string `mapstructure:"BROKER_PASSWORD"`
	BrokerHost           string `mapstructure:"BROKER_HOST"`
	WebServerPort        string `mapstructure:"WEB_SERVER_PORT"`
}

func LoadConfig(env string) (*Conf, error) {
	var cfg *Conf

	if env == "dev" {
		if !fileExists(".env") {
			return nil, fmt.Errorf(".env file not found in dev environment")
		}

		viper.SetConfigFile(".env")
		err := viper.ReadInConfig()
		if err != nil {
			return nil, fmt.Errorf("error reading .env file: %v", err)
		}
		fmt.Println("Loaded .env for dev environment")
	} else {
		fmt.Println("Using system environment variables for prod environment")
	}

	viper.AutomaticEnv()

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}

	return cfg, nil
}

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}
