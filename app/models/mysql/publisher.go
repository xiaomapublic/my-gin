package mysql

import (
	"github.com/jinzhu/gorm"
	"my-gin/app/libraries/mysql"
	"sync"
)

type Publisher struct {
	mysql.Base
	Name string `json:"name"`
	Code string `json:"code"`
	Status int `json:"status"`
	Screenshot_path string `json:"screenshot_path"`
}

func PublisherObj() *gorm.DB {
	var once sync.Once
	var mysqlObj *gorm.DB
	once.Do(func() {
		mysqlObj = mysql.GetORMByName("my_gin").Table("aa_publisher")
	})
	return mysqlObj
}

