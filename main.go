package main

import (
	"context"
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"my-gin/app/cronjobs"
	. "my-gin/app/libraries/config"
	"my-gin/app/libraries/log"
	"my-gin/app/libraries/mongodb"
	"my-gin/app/libraries/mysql"
	"my-gin/app/libraries/redis"
	"my-gin/routers"
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

	//初始化定时任务
	cronjobs.Init()

	//设置cpu最大执行数量
	runtime.GOMAXPROCS(runtime.NumCPU())
	logger.Info("cup核数：", runtime.NumCPU())

	//设置系统模式release为开发模式
	gin.SetMode(DefaultConfig.GetString("mode"))

	//获取gin初始化实例
	router := gin.Default()

	// 性能分析工具
	ginpprof.Wrap(router)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该路由",
		}) 
		return
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该方法",
		})
		return
	})

	routers.RegisterApiRouter(router)

	//gin默认监听端口方式
	//if err := router.Run(DefaultConfig.GetString("server_port")); err != nil {
	//	log.Fatalf("listen: %s\n", err)
	//}

	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome")
	})

	srv := &http.Server{
		Addr:    DefaultConfig.GetString("server_port"),
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
