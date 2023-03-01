package global

import (
	"github.com/wenhao26/toolbox/example/config"
)

// 把viper的解析出来的数据存储,这样每个go文件都已引用global中的配置数据

var (
	Settings config.ExampleConfig
)
