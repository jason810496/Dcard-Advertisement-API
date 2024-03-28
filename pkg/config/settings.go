package config

import (
	"flag"
	"fmt"

	"github.com/jason810496/Dcard-Advertisement-API/pkg/utils"
	"github.com/spf13/viper"
)

// Settings is the configuration instance
var APP_MODE string
var Settings *ServerConfig

func Init() {
	flagSettings()

	viper.SetConfigName(getConfigName())
	viper.AddConfigPath(".env")
	viper.AddConfigPath("../../.env") // for test
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error load config file: %s \n", err))
	}

	Settings = &ServerConfig{}
	if err := viper.Unmarshal(Settings); err != nil {
		panic(fmt.Errorf("fatal error parse config file: %s \n", err))
	}

	utils.PrintJson(Settings)
}

func flagSettings() {
	flag.StringVar(&APP_MODE, "config", "dev", "config file")
	flag.Usage = func() {
		fmt.Println("Usage: ./api -config [mode]")
		fmt.Println("mode: dev, test, prod (default: dev)")
	}
	if !flag.Parsed() {
		flag.Parse()
	}
}

// get config file path by argument
func getConfigName() string {

	switch APP_MODE {
	case "local", "dev", "test", "prod", "staging","kubernetes":
		return APP_MODE
	default:
		panic("Invalid app mode  `" + APP_MODE + "`")
	}
}
