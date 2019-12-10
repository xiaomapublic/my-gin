package defaultExecution

import (
	"fmt"
)

// 默认执行程序
func init() {

	go func() {
		MonitorAdHourMq()
	}()

	go func() {
		MonitorAdHourMqTwo()
	}()

	go func() {
		// 监听键盘输入，根据输入执行指定任务
		for {
			var task string
			fmt.Scanln(&task)
			switch task {
			case "anchorMonitorScript":
				AnchorMonitorScript()
			case "MyGinScript":
				MyGinScript()
			case "nihao":
				fmt.Println("输出nihao")
			case "hello":
				fmt.Println("输出hello")
			case "MyGinElastic":
				MyGinElastic()
			}
		}

	}()

}
