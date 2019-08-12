//配置文件服务类，使用viper库
package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

//默认配置，供全局使用
var DefaultConfig *viper.Viper

//反序列化后的配置，供全局使用
var UnmarshalConfig YamlSetting

var once sync.Once

//单例模式获取默认配置数据
func DefaultConfigInit() {
	//gin框架会优先调用控制器里面init方法，配置文件需要注意
	once.Do(func() {
		DefaultConfig = NewConfig("configs", "config", "yaml")
		UnmarshalConfig = ParseYaml()
	})
}

func NewConfig(filePath string, fileName string, suffix string) *viper.Viper {

	var config *viper.Viper

	config = viper.New() //返回初始化实例
	//DefaultConfig.SetEnvPrefix(cmdRoot)//设置环境变量前缀
	//DefaultConfig.AutomaticEnv()//检查变量设置键
	//replacer := strings.NewReplacer(".", "_")
	//DefaultConfig.SetEnvKeyReplacer(replacer)//用于将环境变量设置到有此功能的键
	config.SetConfigName(fileName)                         //设置文件名
	config.AddConfigPath("$GOPATH/src/my-gin/" + filePath) //设置路径名
	config.SetConfigType(suffix)                           //设置文件后缀
	err := config.ReadInConfig()                           // Find and read the config file
	if err != nil {                                        // Handle errors reading the config file
		fmt.Println(fmt.Errorf("Fatal error when reading %s config file: %s\n", filePath, err))
	}

	return config
}

//反序列化配置文件为结构体
func ParseYaml() YamlSetting {
	var yamlObj YamlSetting
	if err := DefaultConfig.Unmarshal(&yamlObj); err != nil {
		fmt.Printf("err:%s", err)
	}
	return yamlObj
}

type Jwt struct {
	Secret string
	Alg    string
}

type Log struct {
	Path        string
	Max_size    int
	Max_backups int
	Max_age     int
	Compress    bool
}

type Mysql struct {
	Host         string
	User         string
	Passwd       string
	Port         string
	Maxidleconns int
	Maxopenconns int
}

type Redis struct {
	Addr       string
	Pwd        string
	Instance   string
	Max_idle   int
	Max_active int
}

type Mongodb struct {
	Addr       []string
	User       string
	Pwd        string
	Instance   string
	Max_active int
}

type Rabbitmq struct {
	Addr string
	User string
	Pwd  string
}

type YamlSetting struct {
	Mode        string
	Server_port string
	App_name    string
	Template    string
	App_secret  string
	Cookie_name string
	Jwt         Jwt
	Log         Log
	Excel       string
	Mysql       map[string]Mysql
	Redis       map[string][]Redis
	Mongodb     map[string]Mongodb
	Rabbitmq    map[string]Rabbitmq
}
