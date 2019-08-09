//示例代码
package test

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"my-gin/app/libraries/mongodb"
	"my-gin/app/libraries/mysql"
	"my-gin/app/libraries/rabbitmq"
	redisLib "my-gin/app/libraries/redis"
	mongodbMod "my-gin/app/models/mongodb"
	mysqlMod "my-gin/app/models/mysql"
	"my-gin/app/services/test"
	"my-gin/filters/auth"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Api struct {
}

//mysql写入数据
//请求参数示例：{"hour":"2019-09-08 23:00:00","ad_id":"25982059966487040","campaign_id":"25838044971136768","product_id":146,"advertiser_id":103,"request_count":4594,"cpm_count":1076,"cpc_original_count":2,"division_id":3,"status":5}
func (a *Api) MysqlCreate(c *gin.Context) {

	var data mysqlMod.My_gin
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = mysqlMod.My_gin_obj().Create(&data).Error
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "新增成功",
			"data": "",
		})
	}

}

//mysql更新数据
//请求参数示例：{"hour":"2019-09-08 23:00:00","ad_id":"25982059966487040","status":5}
func (a *Api) MysqlUpdate(c *gin.Context) {

	var data mysqlMod.My_gin
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}
	data.Updated_at = time.Now()

	err = mysqlMod.My_gin_obj().Model(&data).Where("hour = ? AND ad_id = ?", data.Hour, data.Ad_id).Update("status", data.Status).Error

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "更新成功",
			"data": "",
		})
	}
}

//mysql删除数据
//请求参数示例：{"hour":"2019-09-08 23:00:00","ad_id":"25982059966487040"}
func (*Api) MysqlDelete(c *gin.Context) {

	var data mysqlMod.My_gin
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = mysqlMod.My_gin_obj().Where("hour = ? AND ad_id = ?", data.Hour, data.Ad_id).Delete(&data).Error

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "删除成功",
			"data": "",
		})
	}

}

//mysql获取全部数据
func (*Api) MysqlGetAll(c *gin.Context) {
	var data []mysqlMod.My_gin
	err := mysqlMod.My_gin_obj().Find(&data).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"data": data,
	})
}

//mysql根据条件查询,调用RABBITMQ
func (*Api) MysqlGetWhere(c *gin.Context) {
	start_time := c.Query("start_time")
	end_time := c.Query("end_time")
	page, _ := strconv.Atoi(c.Query("page"))
	_, adOk := c.GetQuery("ad_id")
	_, caOk := c.GetQuery("campaign_id")
	_, prOk := c.GetQuery("product_id")
	_, advOk := c.GetQuery("advertiser_id")

	var data []mysqlMod.My_gin

	Db := mysql.GetORMByName("my_gin")

	pageTotal := 1000
	offset := (page - 1) * pageTotal

	if adOk == true {
		Db.Raw("SELECT * FROM my_gin WHERE `hour` >= ? AND `hour` <= ? ORDER BY `hour` DESC LIMIT ?,?", start_time, end_time, offset, pageTotal).Scan(&data)
	} else if caOk == true {
		Db.Raw("select `hour`,`campaign_id`, sum(request_count) as request_count, sum(cpm_count) as cpm_count, sum(cpc_original_count) as cpc_original_count from my_gin WHERE `hour` >= ? AND `hour` <= ? GROUP BY `hour`,`campaign_id` ORDER BY `hour` DESC LIMIT ?,?", start_time, end_time, offset, pageTotal).Scan(&data)
	} else if prOk == true {
		Db.Raw("select `hour`,product_id, sum(request_count) as request_count, sum(cpm_count) as cpm_count, sum(cpc_original_count) as cpc_original_count from my_gin WHERE `hour` >= ? AND `hour` <= ? GROUP BY `hour`,product_id ORDER BY `hour` DESC LIMIT ?,?", start_time, end_time, offset, pageTotal).Scan(&data)
	} else if advOk == true {
		Db.Raw("select `hour`,advertiser_id, sum(request_count) as request_count, sum(cpm_count) as cpm_count, sum(cpc_original_count) as cpc_original_count from my_gin WHERE `hour` >= ? AND `hour` <= ? GROUP BY `hour`,advertiser_id ORDER BY `hour` DESC LIMIT ?,?", start_time, end_time, offset, pageTotal).Scan(&data)
	}

	ch := rabbitmq.RabbitSession["my_vhost"]

	//创建交换器
	//err := ch.ExchangeDeclare("", "direct", true, true, false, false, nil)

	//使用默认交换器
	//创建队列
	q, err := ch.QueueDeclare(
		"adHour", // name  有名字！
		true,     // durable  持久性的,如果事前已经声明了该队列，不能重复声明
		false,    // delete when unused
		false,    // exclusive 如果是真，连接一断开，队列删除
		false,    // no-wait
		nil,      // arguments
	)

	//err = ch.QueueBind("adHour", "ad", "st", false, nil)

	rabbitmq.FailOnError(err, "Failed to declare a queue")

	body, _ := json.Marshal(data)

	// 发布
	err = ch.Publish(
		"",     // exchange 默认模式，exchange为空
		q.Name, // routing key 默认模式路由到同名队列，即是task_queue
		false,  // mandatory
		false,
		amqp.Publishing{
			// 持久性的发布，因为队列被声明为持久的，发布消息必须加上这个（可能不用），但消息还是可能会丢，如消息到缓存但MQ挂了来不及持久化。
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})

	rabbitmq.FailOnError(err, "Failed to publish a message")

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"data": data,
	})
}

//redis写入数据
//请求参数示例：{"hour":"2019-09-08 23:00:00","ad_id":"25982059966487041","campaign_id":"25838044971136768","product_id":146,"advertiser_id":103,"request_count":4594,"cpm_count":1076,"cpc_original_count":2,"division_id":3,"status":5}
func (*Api) RedisCreate(c *gin.Context) {
	var data mysqlMod.My_gin
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	//定义redis实例类
	var redisClass redisLib.RedisInstanceClass
	//根据配置获取某个具体redis实例
	redisClass.GetRedigoByName("default", "master")
	//string
	redisClass.Set("set", data.Ad_id)
	redisClass.SetEx("setex", data.Hour, 100)
	//hash
	redisClass.HMSet("hash_"+data.Ad_id, data)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": "",
	})
}

//redis更新数据
//请求参数示例：{"hour":"2019-09-11 23:00:00","ad_id":"25982059966487041","request_count":5005,"cpm_count":4004,"status":100}
func (*Api) RedisUpdate(c *gin.Context) {

	var data mysqlMod.My_gin
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	var redisClass redisLib.RedisInstanceClass
	redisClass.GetRedigoByName("default", "master")
	//string
	redisClass.Set("set", data.Ad_id)
	redisClass.SetEx("setex", data.Hour, 100)
	//hash
	params := make(map[string]interface{})
	params["Hour"] = data.Hour
	params["Request_count"] = data.Request_count
	params["Status"] = data.Status

	redisClass.HMSet("hash_"+data.Ad_id, params)
	redisClass.HSet("hash_"+data.Ad_id, "Request_count", data.Request_count)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": "",
	})
}

//redis删除数据
func (*Api) RedisDelete(c *gin.Context) {

	var data mysqlMod.My_gin
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	var redisClass redisLib.RedisInstanceClass
	redisClass.GetRedigoByName("default", "master")
	redisClass.Del("hash_" + data.Ad_id)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": "",
	})
}

//redis获取指定键
func (*Api) RedisGetWhere(c *gin.Context) {

	ad_id, _ := c.GetQuery("ad_id")

	var redisClass redisLib.RedisInstanceClass
	redisClass.GetRedigoByName("default", "slave")
	//string
	dataString := redisClass.Get("set")
	//hash
	dataHash := redisClass.HMGetAll("hash_" + ad_id)
	dataHashT := redisClass.HGet("hash_"+ad_id, "Product_id")
	c.JSON(http.StatusOK, gin.H{
		"code":       0,
		"msg":        "",
		"dataString": dataString,
		"dataHash":   dataHash,
		"dataHashT":  dataHashT,
	})
}

//mongodb写入数据
//请求参数示例：{"hour":"2019-09-08 23:00:00","ad_id":"25982059966487041","campaign_id":"25838044971136768","product_id":146,"advertiser_id":103,"request_count":4594,"cpm_count":1076,"cpc_original_count":2,"division_id":3,"status":5}
func (*Api) MongodbCreate(c *gin.Context) {
	var data mongodbMod.MyGinData
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	data.Id = mongodb.CreateId()
	data.Created_at = time.Now()

	var conn *mongodbMod.MyGin

	err = conn.Mongodb().Insert(data)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "ok",
			"data": "",
		})
	}
}

//mongodb更新数据
//请求参数示例：{"hour":"2019-09-11 23:00:00","ad_id":"25982059966487041","request_count":5005,"cpm_count":4004,"status":100}
func (*Api) MongodbUpdate(c *gin.Context) {
	var data mongodbMod.MyGinData
	var params = make(map[string]interface{})
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	var conn *mongodbMod.MyGin
	params["request_count"] = data.Request_count
	params["cpm_count"] = data.Cpm_count
	params["status"] = data.Status

	datass := bson.M{"$set": params} //$set关键字表示只更新指定字段
	_, err = conn.Mongodb().UpdateAll(bson.M{"ad_id": data.Ad_id, "hour": data.Hour}, datass)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "ok",
			"data": "",
		})
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

	err = conn.Mongodb().Remove(bson.M{"_id": data.Id})
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "ok",
			"data": "",
		})
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
	ad_id, _ := c.GetQuery("ad_id")

	err := conn.Mongodb().Find(bson.M{"ad_id": ad_id}).All(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  0,
		"data": data,
	})
}

//登陆获取token
func (*Api) JwtSetLogin(c *gin.Context) {
	var data mysqlMod.User_info
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "请输入账号密码",
		})
		return
	}

	var count int
	mysqlMod.User_info_obj().Where("name = ? AND pwd = ?", data.Name, data.Pwd).Count(&count)

	if count < 1 {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "账号密码错误",
		})
		return
	}

	info := make(map[string]interface{})
	info["name"] = data.Name
	info["pwd"] = data.Pwd

	authDr, _ := c.MustGet("jwt_auth").(*auth.Auth)

	token, _ := (*authDr).Login(c.Request, c.Writer, info).(string)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": gin.H{
			"token": token,
		},
	})
}

//通过token获取用户信息
func (*Api) JwtGetUserInfo(c *gin.Context) {
	authDr, _ := c.MustGet("jwt_auth").(*auth.Auth)

	info := (*authDr).User(c).(map[string]interface{})

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": info,
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

/**
 * 返回大写字符串
 */
func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制

	return strings.ToUpper(md5str1)
}

//返回小写字符串
func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制

	return md5str1
}
