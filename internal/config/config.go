package config

import (
	"fmt"

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

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config: %v", err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}
	return cfg, nil
}
