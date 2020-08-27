package config

import (
	"os"

	"github.com/spf13/viper"
)

// InitConfig reads config file to viper instance
func InitConfig(filename string) error {
	var err error

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		return err
	}

	viper.SetConfigFile(filename)
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
