package main

import (
	"fmt"
	"time"

	"github.com/wenhao26/toolbox/example/global"
	"github.com/wenhao26/toolbox/example/initialize"
)

func main() {
	initialize.InitConfig()
	// 此处用来模拟 example.yaml 配置项值改变时，进行热加载无需重启
	for {
		fmt.Println(global.Settings.Name)
		time.Sleep(2e9)
	}
}
