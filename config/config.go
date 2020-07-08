package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

func InitConfig() {
	configFile := flag.String("conf", "./config/config.yaml", "path of config file")
	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()
	if err != nil {
		errStr := fmt.Sprintf("viper read config is failed, err is %v configFile is %s ", err, *configFile)
		panic(errStr)
	}
}
