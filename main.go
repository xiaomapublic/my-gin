package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"my-gin/app/cronjobs"
	"my-gin/app/services/test"
	. "my-gin/libraries/config"
	"my-gin/libraries/log"
	"my-gin/libraries/mongodb"
	"my-gin/libraries/mysql"
	"my-gin/libraries/rabbitmq"
	"my-gin/libraries/redis"
	routerBase "my-gin/libraries/router"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
)

// 应用主函数入口
func main() {
	//gin框架会优先加载路由，会调用控制器里面init方法，配置文件需要注意
	DefaultConfigInit()

	//初始化日志文件
	logger := log.InitLog("main")

	//初始化数据库,包一级声明的变量来说，它们的生命周期和整个程序的运行周期是一致的
	mysql.Init()
	redis.Init()
	mongodb.Init()
	rabbitmq.Init()

	//初始化定时任务
	cronjobs.Init()

	//初始化队列任务
	go func() {
		test.MonitorAdHourMq()
	}()

	//设置cpu最大执行数量
	runtime.GOMAXPROCS(runtime.NumCPU())
	logger.Info("cup核数：", runtime.NumCPU())

	//设置系统模式release为开发模式
	gin.SetMode(UnmarshalConfig.Mode)

	//获取gin初始化实例
	router := routerBase.InitRouter()

	//gin默认监听端口方式
	//if err := router.Run(UnmarshalConfig.Server_port); err != nil {
	//	log.Fatalf("listen: %s\n", err)
	//}

	srv := &http.Server{
		Addr:    UnmarshalConfig.Server_port,
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen: ", err.Error())
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown: ", err.Error())
	}
	logger.Info("Server exiting")

	pid := fmt.Sprintf("%d", os.Getpid())
	_, openErr := os.OpenFile("pid", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if openErr == nil {
		ioutil.WriteFile("pid", []byte(pid), 0)
	}
}
