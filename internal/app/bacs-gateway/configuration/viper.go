package configuration

import (
	"github.com/spf13/viper"
	"strings"
)

func initViper() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigFile("configs/application.yml")
	err := viper.ReadInConfig()
	viper.MergeInConfig()

	if err != nil {
		panic(err)
	}
}
