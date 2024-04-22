package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	_ "github.com/spf13/viper/remote"

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

func initViperRemote() {
	// etcdctl --endpoints=127.0.0.1:12379 put /webook "$(dev.yaml)"
	viper.SetConfigType("ymal")
	err := viper.AddRemoteProvider("etcd3",
		// 通过webook和其他使用etcd的区别出来
		"127.0.0.1:12379", "/webool")
	if err != nil {
		panic(err)
	}

	err = viper.WatchRemoteConfig()
	if err != nil {
		panic(err)
	}
	// viper的OnConfigChange 无法左右远程监听修改
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func InitViperV1() {
	cfile := pflag.String("config", "config/config.yaml", "指定配置文件路径")
	pflag.Parse()

	viper.SetConfigFile(*cfile)

	// 实时监听配置变更
	viper.WatchConfig()

	// 只能告诉文件变了，但未告诉文件哪些变了
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
	})

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
