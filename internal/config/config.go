package config

import "github.com/spf13/viper"

type Config struct {
	DBHost string 
	DBPort int 
	DBUser string 
	DBPassword string 
	DBName string
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("DB_PORT", 5432)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	_ = viper.ReadInConfig() // ignore error if file missing
	viper.AutomaticEnv()

	cfg := &Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetInt("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
	}

	return cfg, nil
}