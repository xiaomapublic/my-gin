package main

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"my-gin/app/libraries/config"
	"my-gin/filters"
	"my-gin/filters/auth"
	routeRegister "my-gin/routers"
	"net/http"
)

func initRouter() *gin.Engine {
	router := gin.New()

	// html模板
	router.LoadHTMLGlob(config.UnmarshalConfig.Template + "/*")

	//性能分析
	ginpprof.Wrap(router)

	router.Use(gin.Logger())

	// 错误处理
	router.Use(handleErrors())
	// 全局session
	router.Use(filters.RegisterSession())
	// 全局cache
	router.Use(filters.RegisterCache())
	// 全局auth cookie
	router.Use(auth.RegisterGlobalAuthDriver("cookie", "web_auth"))
	// 全局auth jwt
	router.Use(auth.RegisterGlobalAuthDriver("jwt", "jwt_auth"))

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

	routeRegister.RegisterApiRouter(router)

	return router
}
