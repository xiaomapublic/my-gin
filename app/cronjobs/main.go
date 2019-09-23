//定时任务主文件，负责配置定时任务运行时间
package cronjobs

import (
	"github.com/robfig/cron"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		runTask()
	}()

	//阻塞主线程
	//var run chan bool
	//<-run
}

func runTask() {
	c := cron.New()
	c.AddFunc("1 0 * * * *", Spider)
	c.Start()
}
