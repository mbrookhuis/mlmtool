package readconfig

import (
	"fmt"
	"os"

	"github.com/go-viper/mapstructure/v2"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ReadConfig
// Read the given configuration file.
// param: cfgFile
// return:
//
//	error: when there is an error this will be returned, otherwise nil
//	model.Config: configuration file data
func ReadConfig(cfgFile string, appConfig interface{}) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		cfgFile = "./config.yaml"
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	fmt.Println("Using config file:", cfgFile)
	if _, err := os.Stat(cfgFile); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Warning: No config file found. Using defaults and environment variables. %w", err)
		return fmt.Errorf("no config file found")
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Warning: No config file found. Using defaults and environment variables.")
		}
	}
	err := viper.Unmarshal(&appConfig, func(dc *mapstructure.DecoderConfig) {})
	if err != nil {
		return fmt.Errorf("unable to decode config into struct: %w", err)
	}
	return nil
}
