package initialize

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/wenhao26/toolbox/example/config"
	"github.com/wenhao26/toolbox/example/global"
)

func InitConfig() {
	v := viper.New()
	//v.SetConfigName("example") // 配置文件的名称(不带扩展名)
	//v.SetConfigType("yaml") // 配置文件扩展名称
	//v.AddConfigPath("./toolbox") // 查找配置文件的路径
	v.SetConfigFile("F:\\toolbox\\example.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	exampleConfig := config.ExampleConfig{}
	if err := v.Unmarshal(&exampleConfig); err != nil {
		log.Println(err.Error())
	}
	global.Settings = exampleConfig

	// 监听配置变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		_ = v.Unmarshal(&exampleConfig)
		global.Settings = exampleConfig
		fmt.Printf("Config file:%s Op:%s\n", e.Name, e.Op)
	})
}
