//定时任务主文件，负责配置定时任务运行时间
package cronjobs

import (
	"github.com/robfig/cron"
	"runtime"
)

func Init() {
	main()
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	go runTask()

	//阻塞主线程
	//var run chan bool
	//<-run
}

func runTask() {
	c := cron.New()
	//c.AddFunc("5,30 * * * * *", Spider)
	//c.AddFunc("15 00 19 * * *", My_gin_script)
	c.Start()
}
