//excel使用
package cronjobs

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"my-gin/app/libraries/log"
	"time"
)

func Spider() {

	nowTime := time.Now()

	f := excelize.NewFile()
	index := f.NewSheet("Sheet2")
	f.SetCellValue("Sheet1", "A1", "100")

	f.SetActiveSheet(index)

	err := f.SaveAs("./static/excel/" + nowTime.Format("2006-01-02") + ".xlsx")

	if err != nil {
		fmt.Println(err)
		log.InitLog("spider").Error("创建excel错误")
	}
}
