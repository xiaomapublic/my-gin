//excel使用
package cronjobs

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"my-gin/libraries/config"
	"my-gin/libraries/log"
	"os"
	"time"
)

func Spider() {

	nowTime := time.Now()

	f := excelize.NewFile()
	index := f.NewSheet("Sheet2")
	f.SetCellValue("Sheet1", "A1", "100")

	f.SetActiveSheet(index)

	path := config.DefaultConfig.GetString("excel")

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(path, os.ModePerm)
		}
	}

	err = f.SaveAs(path + nowTime.Format("2006-01-02") + ".xlsx")

	if err != nil {
		fmt.Println(err)
		log.InitLog("spider").Errorf("创建excel错误", err.Error())
	}
}
