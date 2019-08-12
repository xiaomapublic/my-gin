package router

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	routeRegister "my-gin/configs"
	"my-gin/libraries/config"
	"my-gin/libraries/filters"
	"my-gin/libraries/filters/auth"
	"my-gin/libraries/handle"
	"net/http"
)

func InitRouter() *gin.Engine {
	router := gin.New()

	// html模板
	router.LoadHTMLGlob(config.UnmarshalConfig.Template + "/*")

	//性能分析
	ginpprof.Wrap(router)

	router.Use(gin.Logger())

	// 错误处理
	router.Use(handle.HandleErrors())
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
