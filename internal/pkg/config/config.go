package config

import (
	"os"
	l "telegram-bot-feedback/internal/pkg/logger"

	"github.com/spf13/viper"
)

// GetConfig returns configuration
func GetConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.SetConfigType("json")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			v, err = createConfig(v)
			if err != nil {
				return nil, l.Err(err)
			}
		} else {
			return nil, l.Err(err)
		}
	}
	return v, nil
}

// createConfig creates config
func createConfig(v *viper.Viper) (*viper.Viper, error) {
	file, _ := os.Create("config.json")
	file.Close()
	v.Set("host", "")
	v.Set("token", "")
	v.Set("offset", 0)
	if err := v.WriteConfig(); err != nil {
		return nil, l.Err(err)
	}
	return v, nil
}
