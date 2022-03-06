package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	LOCAL_ADDR         string `mapstructure:"LOCALHOSTt_Addr"`
	SECRET_KEY         string `mapstructure:"SECRET_key"`
	SECRET_REFRESH_KEY string `mapstructure:"SECRET_REFRESH_key"`
	MONGODB_ADDR       string `mapstructure:"MONGODB_uri"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
