//mysql服务类，使用gorm库
package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "my-gin/libraries/config" //  . 操作是省略前面的包名; _ 匿名导入，不适用包内任何方法，但会触发init()方法
	"my-gin/libraries/log"
	"time"
)

var Gorm map[string]*gorm.DB

func init() {
	Gorm = make(map[string]*gorm.DB)
	NewDB()
}

type base interface {
	List()
	Count()
}

// Base 为数据库核心model，其他model均内嵌此model。
type Base struct {
	base
	Id         uint32    `gorm:"primary_key" json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"Updated_at"`
}

// 初始化
func NewDB() {
	var orm *gorm.DB
	var err error

	//viper读出的配置数据会强制将键中的大写转换为小写
	dbConfig := UnmarshalConfig.Mysql

	for dbName, dbConfigs := range dbConfig {

		//打开数据库连接
		orm, err = gorm.Open("mysql", dbConfigs.User+":"+dbConfigs.Passwd+"@tcp("+dbConfigs.Host+":"+dbConfigs.Port+")/"+dbName+"?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			log.InitLog("mysql").Errorf("NewDB", "msg", err.Error())
			fmt.Println(err.Error())
			return
		}

		//建表时不将表名自动变更为单词的复数形式
		orm.SingularTable(true)
		//设置最大空闲数
		orm.DB().SetMaxIdleConns(dbConfigs.Maxidleconns)
		//设置最大连接数
		orm.DB().SetMaxOpenConns(dbConfigs.Maxopenconns)
		//设置每个连接的过期时间
		orm.DB().SetConnMaxLifetime(time.Second * 5)

		Gorm[dbName] = orm

	}

}

// 通过名称获取Gorm实例
func GetORMByName(dbName string) *gorm.DB {

	return Gorm[dbName]
}
