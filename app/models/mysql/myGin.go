package mysql

import (
	"github.com/jinzhu/gorm"
	"my-gin/libraries/mysql"
	"time"
)

type MyGin struct {
	Hour               string    `json:"hour"`
	Ad_id              string    `json:"ad_id"`
	Campaign_id        string    `json:"campaign_id"`
	Product_id         int       `json:"product_id"`
	Advertiser_id      int       `json:"advertiser_id"`
	Request_count      int       `json:"request_count"`
	Cpm_count          int       `json:"cpm_count"`
	Cpc_original_count int       `json:"cpc_original_count"`
	Division_id        int       `json:"division_id"`
	Status             int       `json:"status"`
	Created_at         time.Time `json:"created_at"`
	Updated_at         time.Time `json:"updated_at"`
}

func MyGinObj() *gorm.DB {

	var mysqlObj *gorm.DB

	mysqlObj = mysql.GetORMByName("my_gin").Table("my_gin")

	return mysqlObj
}
