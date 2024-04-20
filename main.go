package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	InitViperV1()
	server := InitWebServer()
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func InitViperV1() {
	cfile := pflag.String("config", "config/config.yaml", "指定配置文件路径")
	pflag.Parse()

	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func InitViper() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
