package config

import (
	"os"

	"github.com/spf13/viper"
)

const (
	EnvConfigFile         = "CONF_FILE"
	DefaultConfigLocation = "."
	DefaultConfigFile     = "config"
)

type Config interface {
	Get(key string) string
}

func Load() (Config, error) {
	v := viper.New()

	configFile := os.Getenv(EnvConfigFile)
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName(DefaultConfigFile)
		v.AddConfigPath(DefaultConfigLocation)
	}
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &viperConfig{viper: v}, nil
}

type viperConfig struct {
	viper *viper.Viper
}

func (vc *viperConfig) Get(key string) string {
	if vc.viper.IsSet(key) {
		return vc.viper.GetString(key)
	} else {
		return ""
	}
}
