//excel使用
package cronjobs

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"my-gin/app/libraries/config"
	"my-gin/app/libraries/log"
	"time"
)

func Spider() {

	nowTime := time.Now()

	f := excelize.NewFile()
	index := f.NewSheet("Sheet2")
	f.SetCellValue("Sheet1", "A1", "100")

	f.SetActiveSheet(index)

	config.DefaultConfigInit()
	path := config.DefaultConfig.GetString("excel")

	err := f.SaveAs(path + nowTime.Format("2006-01-02") + ".xlsx")

	if err != nil {
		fmt.Println(err)
		log.InitLog("spider").Errorf("创建excel错误", err.Error())
	}
}
