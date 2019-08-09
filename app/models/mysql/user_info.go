package mysql

import (
	"github.com/jinzhu/gorm"
	"my-gin/app/libraries/mysql"
	"time"
)

type User_info struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Pwd        string    `json:"pwd"`
	Status     int       `json:"status"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func User_info_obj() *gorm.DB {

	var mysqlObj *gorm.DB

	mysqlObj = mysql.GetORMByName("my_gin").Table("user_info")

	return mysqlObj
}
