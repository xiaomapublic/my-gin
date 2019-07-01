package mysql

import (
	"github.com/jinzhu/gorm"
	"my-gin/app/libraries/mysql"
)

type Publisher struct {
	mysql.Base
	Name            string `json:"name"`
	Code            string `json:"code"`
	Status          int    `json:"status"`
	Screenshot_path string `json:"screenshot_path"`
}

func PublisherObj() *gorm.DB {

	var mysqlObj *gorm.DB
	mysqlObj = mysql.GetORMByName("my_gin").Table("aa_publisher")

	return mysqlObj
}
