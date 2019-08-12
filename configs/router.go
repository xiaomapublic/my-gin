package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"my-gin/app/controllers/test"
	_ "my-gin/docs"
	"my-gin/libraries/filters/auth"
)

var testApi *test.Api

func RegisterApiRouter(router *gin.Engine) {
	//swagger文档
	router.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))

	//示例接口路由
	testApiRouter := router.Group("test")
	{

		testApiRouter.POST("/jwt/set", testApi.JwtSetLogin)

		testApiRouter.Use(auth.Middleware("jwt"))
		{

			testApiRouter.POST("/Api/MysqlCreate", testApi.MysqlCreate)
			testApiRouter.POST("/Api/MysqlUpdate", testApi.MysqlUpdate)
			testApiRouter.POST("/Api/MysqlDelete", testApi.MysqlDelete)
			testApiRouter.GET("/Api/MysqlGetAll", testApi.MysqlGetAll)
			testApiRouter.GET("/Api/MysqlGetWhere", testApi.MysqlGetWhere)
			testApiRouter.POST("/Api/RedisCreate", testApi.RedisCreate)
			testApiRouter.POST("/Api/RedisUpdate", testApi.RedisUpdate)
			testApiRouter.POST("/Api/RedisDelete", testApi.RedisDelete)
			testApiRouter.GET("/Api/RedisGetWhere", testApi.RedisGetWhere)
			testApiRouter.POST("/Api/MongodbCreate", testApi.MongodbCreate)
			testApiRouter.POST("/Api/MongodbUpdate", testApi.MongodbUpdate)
			testApiRouter.POST("/Api/MongodbDelete", testApi.MongodbDelete)
			testApiRouter.GET("/Api/MongodbGetAll", testApi.MongodbGetAll)
			testApiRouter.GET("/Api/MongodbGetWhere", testApi.MongodbGetWhere)
			testApiRouter.GET("/Api/RandomNumber", testApi.RandomNumber)
			testApiRouter.GET("/Api/Concurrent", testApi.Concurrent)
			testApiRouter.GET("/jwt/get", testApi.JwtGetUserInfo)
		}

	}
}
