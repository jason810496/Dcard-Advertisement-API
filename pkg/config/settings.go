package config

import (
	"github.com/spf13/viper"
	"fmt"
	"os"
)


// Settings is the configuration instance
var Settings *ServerConfig


func init() {
	viper.SetConfigName(getConfigName())
	viper.AddConfigPath(".env")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error load config file: %s \n", err))
	}

	Settings = &ServerConfig{}
	if err := viper.Unmarshal(Settings); err != nil {
		panic(fmt.Errorf("fatal error parse config file: %s \n", err))
	}

	fmt.Printf("%#v\n", *Settings)
}


// GetConfig returns the configuration instance
// func GetConfig() *ServerConfig {
// 	return Settings
// }

// get config file path by argument
func getConfigName() string {
	if len(os.Args) < 2 {
		return "dev.yaml"
	}

	app_mode := os.Args[1]

	switch app_mode {
	case "dev":
		return "dev.yaml"
	case "test":
		return "test.yaml"
	case "prod":
		return "prod.yaml"
	default:
		panic("Invalid app mode")
	}
}
