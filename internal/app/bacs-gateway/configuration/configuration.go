package configuration

import (
	"github.com/spf13/viper"
	"fmt"
)

func GetOrDefault(key string, defaultValue string) string {
	value := viper.GetString(key)
	if value != "" {
		return value
	} else {
		return defaultValue
	}
}

func Get(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func Unmarshal(key string, v interface{}) {
	sub := viper.Sub("db")
	fmt.Println("Sub:", sub)
	viper.UnmarshalKey(key, v)
}

func Start() {
	initViper()
	initVault()
}
