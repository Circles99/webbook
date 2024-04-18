package main

import "github.com/spf13/viper"

func main() {
	InitViper()
	server := InitWebServer()
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func InitViperV1() {
	viper.SetConfigFile("./config/dev.yaml")
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
