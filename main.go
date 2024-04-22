package main

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	InitViper()
	server := InitWebServer()
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func initViperRemote() {
	// etcdctl --endpoints=127.0.0.1:12379 put /webook "$(dev.yaml)"
	viper.SetConfigType("ymal")
	err := viper.AddRemoteProvider("etcd3",
		// 通过webook和其他使用etcd的区别出来
		"127.0.0.1:12379", "/webool")
	if err != nil {
		panic(err)
	}

	err = viper.ReadRemoteConfig()
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
