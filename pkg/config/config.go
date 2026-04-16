package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// BaseConfig chứa các thông số chung cho mọi Microservice
type BaseConfig struct {
	AppName     string `mapstructure:"APP_NAME"`
	AppEnv      string `mapstructure:"APP_ENV"`
	AppPort     int    `mapstructure:"APP_PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
}

// LoadConfig sử dụng Viper để đọc cấu hình từ .env và Environment Variables
func LoadConfig(path string, configName string, out interface{}) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	// Giúp map APP_PORT từ env thành AppPort trong struct
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Warning: .env file not found, using Environment Variables only")
		} else {
			return err
		}
	}

	return viper.Unmarshal(out)
}
