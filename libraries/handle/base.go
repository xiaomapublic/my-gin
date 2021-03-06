package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"my-gin/libraries/log"
	"net/http"
)

func HandleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//初始化日志文件
				logger := log.InitLog("error")
				logger.Errorf("error", err)

				var (
					errMsg     string
					mysqlError *mysql.MySQLError
					ok         bool
				)
				if errMsg, ok = err.(string); ok {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code": 500,
						"msg":  "system error, " + errMsg,
					})
					return
				} else if mysqlError, ok = err.(*mysql.MySQLError); ok {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code": 500,
						"msg":  "system error, " + mysqlError.Error(),
					})
					return
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"code": 500,
						"msg":  fmt.Sprintf("system error %s", err),
					})
					return
				}
			}
		}()
		c.Next()
	}
}
