//示例代码
package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	redisLib "my-gin/app/libraries/redis"
	mongodbMod "my-gin/app/models/mongodb"
	mysqlMod "my-gin/app/models/mysql"
	"my-gin/app/services/test"
	"net/http"
	"sync"
	"time"
)

type Api struct {

}

//mysql写入数据
//请求参数示例：{"name":"my-gin测试","code":"pd-my-gin","status":1,"screenshot_path":"20190508/piAC9MMZeCcU3aXv1uIu.png"}
func (a *Api) MysqlCreate(c *gin.Context) {

	var data mysqlMod.Publisher
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = mysqlMod.PublisherObj().Create(&data).Error
	if err != nil {
		fmt.Println(err.Error())
	}

}

//mysql更新数据
//请求参数示例：{"id":31,"name":"my-gin测试更新"}
func (a *Api) MysqlUpdate(c *gin.Context) {

	var data mysqlMod.Publisher
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}
	data.Updated_at = time.Now()

	err = mysqlMod.PublisherObj().Model(&data).Update(data).Error

	if err != nil {
		fmt.Println(err.Error())
	}
}

//mysql删除数据
//请求参数示例：{"id":36}
func (*Api) MysqlDelete(c *gin.Context) {

	var data mysqlMod.Publisher
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = mysqlMod.PublisherObj().Delete(&data).Error

	if err != nil {
		fmt.Println(err.Error())
	}

}

//mysql获取全部数据
func (*Api) MysqlGetAll(c *gin.Context) {
	var data []mysqlMod.Publisher
	err := mysqlMod.PublisherObj().Find(&data).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"data": data,
	})
}

//mysql根据条件查询
func (*Api) MysqlGetWhere(c *gin.Context) {
	var data []mysqlMod.Publisher

	name,_ := c.GetQuery("name")
	status,_ := c.GetQuery("status")

	err := mysqlMod.PublisherObj().Where("name = ? AND status = ?", name, status).Find(&data).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"data": data,
	})
}

//redis写入数据
func (*Api) RedisCreate(c *gin.Context) {
	//定义redis实例类
	var redisClass redisLib.RedisInstanceClass
	//根据配置获取某个具体redis实例
	redisClass.GetRedigoByName("default", "master")
	//string
	redisClass.Set("set", "yes")
	redisClass.SetEx("setex", "yes", 100)
	//hash
	redisClass.HMSet("hash", map[string]string{"name":"maya", "old":"18"})
	redisClass.HSet("hashT", "name", "xiaoma")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": "",
	})
}

//redis更新数据
func (*Api) RedisUpdate(c *gin.Context) {
	var redisClass redisLib.RedisInstanceClass
	redisClass.GetRedigoByName("default", "master")
	//string
	redisClass.Set("set", "no")
	redisClass.SetEx("setex", "no", 500)
	//hash
	redisClass.HMSet("hash", map[string]string{"name":"xiaoyang"})
	redisClass.HSet("hashT", "name", "xiaoyang")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": "",
	})
}

//redis删除数据
func (*Api) RedisDelete(c *gin.Context) {
	var redisClass redisLib.RedisInstanceClass
	redisClass.GetRedigoByName("default", "master")
	redisClass.Del("status")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": "",
	})
}

//redis获取指定键
func (*Api) RedisGetWhere(c *gin.Context) {
	var redisClass redisLib.RedisInstanceClass
	redisClass.GetRedigoByName("default", "master")
	//string
	dataString := redisClass.Get("set")
	//hash
	dataHash := redisClass.HMGetAll("hash")
	dataHashT := redisClass.HGet("hash", "name")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"dataString": dataString,
		"dataHash": dataHash,
		"dataHashT": dataHashT,
	})
}

//mongodb写入数据
func (*Api) MongodbCreate(c *gin.Context) {
	var data mongodbMod.MyGinData
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	data.Id = bson.NewObjectId().Hex()
	data.Created_at = time.Now()

	var conn *mongodbMod.MyGin

	err = conn.Mongodb().Insert(data)
	if err != nil {
		fmt.Println(err.Error())
	}
}

//mongodb更新数据
func (*Api) MongodbUpdate(c *gin.Context) {
	var data mongodbMod.MyGinData
	var params = make(map[string]interface{})
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}
	var conn *mongodbMod.MyGin
	params["name"] = data.Name
	datass := bson.M{"$set": params}//$set关键字表示只更新指定字段
	err = conn.Mongodb().UpdateId(data.Id, datass)
	if err != nil {
		fmt.Println(err.Error())
	}
}

//mongodb删除数据
func (*Api) MongodbDelete(c *gin.Context) {
	var data mongodbMod.MyGinData
	var conn *mongodbMod.MyGin

	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = conn.Mongodb().Remove(bson.M{"_id":data.Id})
	if err != nil {
		fmt.Println(err.Error())
	}
}

//mongodb获取全部数据
func (*Api) MongodbGetAll(c *gin.Context) {
	var conn *mongodbMod.MyGin
	var data []mongodbMod.MyGinData

	err := conn.Mongodb().Find(nil).All(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  0,
		"data": data,
	})
}

//mongodb根据条件查询
func (*Api) MongodbGetWhere(c *gin.Context) {
	var conn *mongodbMod.MyGin
	var data []mongodbMod.MyGinData
	name, _ := c.GetQuery("name")

	err := conn.Mongodb().Find(bson.M{"name":name}).All(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  0,
		"data": data,
	})
}

//生成随机数
func (*Api) RandomNumber(c *gin.Context) {
	randObj := rand.New(rand.NewSource(time.Now().UnixNano()))
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "随机数",
		"data": randObj.Intn(99) + 1,
	})
}

// Go并发编程基础实例
func (*Api) Concurrent(c *gin.Context) {
	people := []string{"Anna", "Bob", "Cody", "Dave", "Eva"}
	match := make(chan string, 1) // 为一个未匹配的发送操作提供空间
	wg := new(sync.WaitGroup)
	wg.Add(len(people))
	for _, name := range people {
		go test.Seek(name, match, wg)
	}
	wg.Wait()
	select {
	case name := <-match:
		fmt.Printf("No one received %s’s message.\n", name)
	default:
		// 没有待处理的发送操作
	}
}
