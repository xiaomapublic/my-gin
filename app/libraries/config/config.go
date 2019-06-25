//配置文件服务类，使用viper库
package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

//默认配置，供全局使用
var DefaultConfig *viper.Viper

var once sync.Once

//单例模式获取默认配置数据
func DefaultConfigInit() {
	//gin框架会优先调用控制器里面init方法，配置文件需要注意
	once.Do(func() {
		DefaultConfig = NewConfig("configs", "config", "yaml")
	})
}

func NewConfig(filePath string, fileName string, suffix string) *viper.Viper {

	var config *viper.Viper

	config = viper.New() //返回初始化实例
	//DefaultConfig.SetEnvPrefix(cmdRoot)//设置环境变量前缀
	//DefaultConfig.AutomaticEnv()//检查变量设置键
	//replacer := strings.NewReplacer(".", "_")
	//DefaultConfig.SetEnvKeyReplacer(replacer)//用于将环境变量设置到有此功能的键
	config.SetConfigName(fileName) //设置文件名
	config.AddConfigPath(filePath) //设置路径名
	config.SetConfigType(suffix)   //设置文件后缀

	err := config.ReadInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		fmt.Println(fmt.Errorf("Fatal error when reading %s config file: %s\n", filePath, err))
	}

	return config
}
