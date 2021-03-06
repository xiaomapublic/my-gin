package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	_ "my-gin/app/cronjobs"
	_ "my-gin/app/services/defaultExecution"
	"my-gin/libraries/config"
	"my-gin/libraries/log"
	_ "my-gin/libraries/mongodb"
	_ "my-gin/libraries/mysql"
	_ "my-gin/libraries/redis"
	routerBase "my-gin/libraries/router"

	"my-gin/app/controllers/test"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	port     = flag.String("port", config.UnmarshalConfig.Server_port, "http port")
	grpcPort = flag.String("grpc-port", config.UnmarshalConfig.Grpc_port, "grpc port")
	mode     = flag.String("mode", config.UnmarshalConfig.Mode, "app mod")
)

// 应用主函数入口
func main() {
	// parse flag
	flag.Parse()
	// 设置系统模式模式
	gin.SetMode(*mode)

	// 设置cpu最大执行数量，go1.8以后的版本不用设置
	// runtime.GOMAXPROCS(runtime.NumCPU())

	logger := log.InitLog("main")
	logger.Info("cup核数：", runtime.NumCPU())
	// 获取gin初始化实例
	router := routerBase.InitRouter()

	// grpc监听
	go func() {
		lis, _ := net.Listen("tcp", ":"+*grpcPort)
		grpcServer := grpc.NewServer()
		grpcApiServer := new(test.GrpcApiServer)
		test.RegisterTestApiServer(grpcServer, grpcApiServer)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("listen-grpc: ", err.Error())
		}
	}()

	// gin默认监听端口方式
	// if err := router.Run(UnmarshalConfig.Server_port); err != nil {
	//	log.Fatalf("listen: %s\n", err)
	// }

	srv := &http.Server{
		Addr:    ":" + *port,
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
