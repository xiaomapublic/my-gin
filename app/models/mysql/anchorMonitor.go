package mysql

import (
	"github.com/jinzhu/gorm"
	"my-gin/libraries/mysql"
)

type AnchorMonitor struct {
	Date             string  `json:"date"`
	Agent_id         int     `json:"agent_id"`
	Agent_channel_id int     `json:"agent_channel_id"`
	Cpm_count        int     `json:"cpm_count"`
	Ecpm             float32 `json:"ecpm"`
	Income           float32 `json:"income"`
	Status           int     `json:"status"`
	Created_at       string  `json:"created_at"`
	Updated_at       string  `json:"updated_at"`
}

func AnchorMonitorObj() *gorm.DB {

	var mysqlObj *gorm.DB

	mysqlObj = mysql.GetORMByName("my_gin").Table("y_anchor_monitor")

	return mysqlObj
}
