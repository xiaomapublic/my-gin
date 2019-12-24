// 示例代码
package test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	mongodbMod "my-gin/app/models/mongodb"
	mysqlMod "my-gin/app/models/mysql"
	"my-gin/app/services/test"
	"my-gin/libraries/config"
	"my-gin/libraries/elastic"
	"my-gin/libraries/filters/auth"
	"my-gin/libraries/mongodb"
	"my-gin/libraries/mysql"
	"my-gin/libraries/rabbitmq"
	redisLib "my-gin/libraries/redis"
	"my-gin/libraries/util"

	"github.com/gin-gonic/gin"
	"github.com/jianfengye/collection"
	elastic2 "github.com/olivere/elastic/v7"
	"github.com/streadway/amqp"
	"github.com/syyongx/php2go"
	"github.com/uniplaces/carbon"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2/bson"
)

type Api struct {
}

// mysql写入数据
// 请求参数示例：{"hour":"2019-09-08 23:00:00","ad_id":"25982059966487040","campaign_id":"25838044971136768","product_id":146,"advertiser_id":103,"request_count":4594,"cpm_count":1076,"cpc_original_count":2,"division_id":3,"status":5}
func (a *Api) MysqlCreate(c *gin.Context) {

	var data mysqlMod.MyGin
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = mysqlMod.MyGinObj().Create(&data).Error
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "新增成功",
			"data": "",
		})
	}

}

// mysql更新数据
// 请求参数示例：{"hour":"2019-09-08 23:00:00","ad_id":"25982059966487040","status":5}
func (a *Api) MysqlUpdate(c *gin.Context) {

	var data mysqlMod.MyGin
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}
	data.Updated_at = time.Now()

	err = mysqlMod.MyGinObj().Model(&data).Where("hour = ? AND ad_id = ?", data.Hour, data.Ad_id).Update("status", data.Status).Error

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "更新成功",
			"data": "",
		})
	}
}

// mysql删除数据
// 请求参数示例：{"hour":"2019-09-08 23:00:00","ad_id":"25982059966487040"}
func (*Api) MysqlDelete(c *gin.Context) {

	var data mysqlMod.MyGin
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = mysqlMod.MyGinObj().Where("hour = ? AND ad_id = ?", data.Hour, data.Ad_id).Delete(&data).Error

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "删除成功",
			"data": "",
		})
	}

}

// mysql获取全部数据
func (*Api) MysqlGetAll(c *gin.Context) {
	var data []mysqlMod.MyGin
	err := mysqlMod.MyGinObj().Limit(10).Find(&data).Error
	if err != nil {
		fmt.Println(err.Error())
	}

	// 时间操作
	nowTime := carbon.Now().AddHour().String()

	// 集合操作
	arr := collection.NewObjCollection(data)

	// 排序
	arr.SortBy("Request_count")

	// 取出一列值
	arrPluck := arr.Pluck("Request_count")
	arrPluck.DD()

	// 根据条件筛选数据
	arrFiter := arr.Filter(func(item interface{}, key int) bool {
		val := item.(mysqlMod.MyGin)
		return val.Request_count > 5000
	})

	// each遍历数组
	request_count := 0
	arrFiter.Each(func(item interface{}, key int) {
		v := item.(mysqlMod.MyGin)
		request_count += v.Request_count
	})

	// 格式化输出
	type Respon struct {
		Hour          string
		Ad_id         string
		Status        string
		Request_count int
	}
	arrEach := collection.NewObjCollection(make([]Respon, 0))
	arrFiter.Each(func(item interface{}, key int) {
		var newMap Respon
		v := item.(mysqlMod.MyGin)

		newMap.Hour = carbon.NewCarbon(v.Hour).String()
		newMap.Ad_id = v.Ad_id
		if v.Status == 5 {
			newMap.Status = "启动"
		} else {
			newMap.Status = "停用"
		}
		newMap.Request_count = v.Request_count
		arrEach.Append(newMap)
	})

	// map重新创建切片
	arrMap := arrEach.Map(func(item interface{}, key int) interface{} {
		v := item.(Respon)
		return v.Request_count
	})

	// reduce聚合计算
	arrReduce := arrMap.Reduce(func(carry collection.IMix, item collection.IMix) collection.IMix {
		carryInt, _ := carry.ToInt()
		itemInt, _ := item.ToInt()
		return collection.NewMix(carryInt + itemInt)
	})

	// 类php插件
	php2go.MbStrlen("nihao你好")

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  nowTime,
		"data": arrReduce,
	})
}

// mysql根据条件查询,调用RABBITMQ
func (*Api) MysqlGetWhere(c *gin.Context) {
	start_time := c.DefaultPostForm("start_time", "2019-07-07 00:00:00")
	end_time := c.DefaultPostForm("end_time", "2019-07-07 00:00:00")
	page, _ := strconv.Atoi(c.DefaultPostForm("page", "1"))
	_, adOk := c.GetPostForm("ad_id")
	_, caOk := c.GetPostForm("campaign_id")
	_, prOk := c.GetPostForm("product_id")
	_, advOk := c.GetPostForm("advertiser_id")

	// fmt.Println(start_time, end_time, page, adOk, caOk, prOk, advOk)

	var data []mysqlMod.MyGin

	Db := mysql.GetORMByName("my_gin")

	pageTotal := 60000
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
	fmt.Printf("%+v", data)
	ch := rabbitmq.Init("my_vhost")
	defer ch.Close()
	// 创建交换器
	// err := ch.ExchangeDeclare("", "direct", true, true, false, false, nil)

	// 使用默认交换器
	// 创建队列
	q, err := ch.QueueDeclare(
		"adHour", // name  有名字！
		true,     // durable  持久性的,如果事前已经声明了该队列，不能重复声明
		false,    // delete when unused
		false,    // exclusive 如果是真，连接一断开，队列删除
		false,    // no-wait
		nil,      // arguments
	)

	// err = ch.QueueBind("adHour", "ad", "st", false, nil)

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

// redis写入数据
// 请求参数示例：{"hour":"2019-09-08 23:00:00","ad_id":"25982059966487041","campaign_id":"25838044971136768","product_id":146,"advertiser_id":103,"request_count":4594,"cpm_count":1076,"cpc_original_count":2,"division_id":3,"status":5}
func (*Api) RedisCreate(c *gin.Context) {
	var data mysqlMod.MyGin
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 定义redis实例类
	var redisClass redisLib.RedisInstanceClass
	// 根据配置获取某个具体redis实例
	redisClass.GetRedigoByName("default")
	// string
	redisClass.Set("set", data.Ad_id)
	redisClass.SetEx("setex", data.Hour, 100)
	// hash
	redisClass.HMSet("hash_"+data.Ad_id, data)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": "",
	})
}

// redis更新数据
// 请求参数示例：{"hour":"2019-09-11 23:00:00","ad_id":"25982059966487041","request_count":5005,"cpm_count":4004,"status":100}
func (*Api) RedisUpdate(c *gin.Context) {

	var data mysqlMod.MyGin
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	var redisClass redisLib.RedisInstanceClass
	redisClass.GetRedigoByName("default")
	// string
	redisClass.Set("set", data.Ad_id)
	redisClass.SetEx("setex", data.Hour, 100)
	// hash
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

// redis删除数据
func (*Api) RedisDelete(c *gin.Context) {

	var data mysqlMod.MyGin
	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	var redisClass redisLib.RedisInstanceClass
	redisClass.GetRedigoByName("default")
	redisClass.Del("hash_" + data.Ad_id)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": "",
	})
}

// redis获取指定键
func (*Api) RedisGetWhere(c *gin.Context) {

	ad_id, _ := c.GetQuery("ad_id")

	var redisClass redisLib.RedisInstanceClass
	redisClass.GetRedigoByName("default")
	// string
	dataString := redisClass.Get("set")
	// hash
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

// mongodb写入数据
// 请求参数示例：{"hour":"2019-09-08 23:00:00","ad_id":"25982059966487041","campaign_id":"25838044971136768","product_id":146,"advertiser_id":103,"request_count":4594,"cpm_count":1076,"cpc_original_count":2,"division_id":3,"status":5}
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

// mongodb更新数据
// 请求参数示例：{"hour":"2019-09-11 23:00:00","ad_id":"25982059966487041","request_count":5005,"cpm_count":4004,"status":100}
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

	datass := bson.M{"$set": params} // $set关键字表示只更新指定字段
	_, err = conn.Mongodb().UpdateAll(bson.M{"ad_id": data.Ad_id}, datass)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "ok",
			"data": "",
		})
	}
}

// mongodb删除数据
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

// mongodb获取全部数据
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

// mongodb根据条件查询
func (*Api) MongodbGetWhere(c *gin.Context) {
	var conn *mongodbMod.MyGin
	var data []mongodbMod.MyGinData
	id, _ := c.GetQuery("id")

	err := conn.Mongodb().Find(bson.M{"_id": id}).All(&data)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  0,
		"data": data,
	})
}

// 登陆获取token
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

	// 获取全局注册的验证驱动程序
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

// 通过token获取用户信息
func (*Api) JwtGetUserInfo(c *gin.Context) {
	authDr, _ := c.MustGet("jwt_auth").(*auth.Auth)

	info := (*authDr).User(c).(map[string]interface{})

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": info,
	})
}

// 生成随机数
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

// mysql大数据量测试
func (*Api) BigDataGet(c *gin.Context) {

	stimeQuery := c.DefaultQuery("stime", carbon.Now().DateString())
	etimeQuery := c.DefaultQuery("etime", carbon.Now().DateString())
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "15"))

	stime, _ := carbon.Parse(carbon.DateFormat, stimeQuery, "Asia/Shanghai")
	etime, _ := carbon.Parse(carbon.DateFormat, etimeQuery, "Asia/Shanghai")
	shour := stime.StartOfDay().String()
	ehour := etime.EndOfDay().String()
	mongoShour := stime.StartOfDay().Local()
	mongoEhour := etime.EndOfDay().Local()
	var (
		err        error
		t1         time.Time
		t2         time.Time
		mysqlCount int
		tidbCount  int
		mongoCount int
		mysqlData  []mysqlMod.MyGin
		tidbData   []mysqlMod.MyGin
		mongoData  []mongodbMod.MyGinAdData
		monConn    *mongodbMod.MyGinAd
	)

	t1 = time.Now()
	if err = mysqlMod.MyGinObj().Where("1=1 AND status=5 AND hour >= ? AND hour <= ?", shour, ehour).Count(&mysqlCount).Error; err != nil {
		fmt.Println(err.Error())
		return
	}
	t2 = time.Now()
	fmt.Println("mysqlCount运行时间：", t2.Sub(t1))

	t1 = time.Now()
	if mongoCount, err = monConn.Mongodb().Find(bson.M{"status": 5, "hour": bson.M{"$gte": mongoShour, "$lte": mongoEhour}}).Count(); err != nil {
		fmt.Println(err.Error())
		return
	}
	t2 = time.Now()
	fmt.Println("mongoCount运行时间：", t2.Sub(t1))

	// t1 = time.Now()
	// if err = mysqlMod.MyGinTidbObj().Where("1=1 AND status=5 AND hour >= ? AND hour <= ?", shour, ehour).Count(&tidbCount).Error; err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// t2 = time.Now()
	// fmt.Println("tidbCount运行时间：", t2.Sub(t1))

	t1 = time.Now()
	if err = mysqlMod.MyGinObj().Where("1=1 AND status=5 AND hour >= ? AND hour <= ?", shour, ehour).Offset(offset).Limit(limit).Order("id").Find(&mysqlData).Error; err != nil {
		fmt.Println(err.Error())
		return
	}
	t2 = time.Now()
	fmt.Println("mysqlData运行时间：", t2.Sub(t1))

	// t1 := time.Now()
	// _ = mysqlMod.MyGinTidbObj().Where("1=1 AND status=5 AND hour >= ? AND hour <= ?", shour, ehour).Offset(offset).Limit(limit).Find(&tidbData).Error
	// t2 := time.Now()
	// fmt.Println("tidbData运行时间：", t2.Sub(t1))
	t1 = time.Now()
	err = monConn.Mongodb().Find(bson.M{"status": 5, "hour": bson.M{"$gte": mongoShour, "$lte": mongoEhour}}).Skip(offset).Limit(limit).All(&mongoData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t2 = time.Now()
	fmt.Println("mongoData运行时间：", t2.Sub(t1))

	c.JSON(http.StatusOK, gin.H{
		"code":       0,
		"mysqlCount": mysqlCount,
		"tidbCount":  tidbCount,
		"mongoCount": mongoCount,
		"mysqlData":  mysqlData,
		"tidbData":   tidbData,
		"mongoData":  mongoData,
	})
}

/**
 * TopK堆排序
 */
func (*Api) TopK(c *gin.Context) {
	countQuery := c.DefaultQuery("count", "100")
	totalQuery := c.DefaultQuery("total", "100")
	op := c.DefaultQuery("op", "asc")

	count, _ := strconv.Atoi(countQuery)
	total, _ := strconv.Atoi(totalQuery)

	type Respon struct {
		// Hour             time.Time `json:"hour"`
		// AdId             string    `json:"ad_id"`
		RequestCount int `json:"request_count" bson:"request_count"`
		// CpmCount         int       `json:"cpm_count"`
		// CpcOriginalCount int       `json:"cpc_original_count"`
	}

	var (
		err        error
		sort       []int
		resultAsc  []int
		resultDesc []int

		page    = 1
		limit   = 20000
		monConn *mongodbMod.MyGinAd
	)

	if total < limit {
		limit = total
	}

	t1 := time.Now()
	for {
		offset := (page - 1) * limit
		var mysqlData []Respon
		// if err = mysqlMod.MyGinObj().Select("request_count").Offset(offset).Limit(limit).Find(&mysqlData).Error; err != nil {
		// 	fmt.Println(err.Error())
		// }
		if err = monConn.Mongodb().Find(bson.M{}).Select(bson.M{"request_count": 1}).Skip(offset).Limit(limit).All(&mysqlData); err != nil {
			fmt.Println(err.Error())
		}
		collection.NewObjCollection(mysqlData).Each(func(item interface{}, key int) {
			v := item.(Respon)
			sort = append(sort, v.RequestCount)
		})

		if len(sort) >= total {
			sort = sort[:total:total]
			break
		}

		if len(mysqlData) < limit {
			break
		}
		page++
	}
	fmt.Printf("sort切片地址：%p，数据总量：%d\n", sort, len(sort))
	t2 := time.Now()
	fmt.Println("readDB运行时间：", t2.Sub(t1), runtime.NumGoroutine())

	if strings.ToUpper(op) == "DESC" {
		maxNumber := util.GetMaxNumber(count, sort)
		fmt.Printf("最大值堆地址：%p\n", maxNumber)
		t3 := time.Now()
		resultAsc = util.SmallHeapAsc(maxNumber)
		fmt.Printf("正序地址：%p\n", resultAsc)
		t4 := time.Now()
		resultDesc = util.SmallHeapDesc(maxNumber)
		fmt.Printf("倒序地址：%p\n", resultDesc)
		t5 := time.Now()
		fmt.Println("Asc运行时间：", t4.Sub(t3))
		fmt.Println("Desc运行时间：", t5.Sub(t4))
	} else {
		minNumber := util.GetMinNumber(count, sort)
		fmt.Printf("最小值堆地址：%p\n", minNumber)
		t3 := time.Now()
		resultAsc = util.LargeHeapAsc(minNumber)
		fmt.Printf("正序地址：%p\n", resultAsc)
		t4 := time.Now()
		resultDesc = util.LargeHeapDesc(minNumber)
		fmt.Printf("倒序地址：%p\n", resultDesc)
		t5 := time.Now()
		fmt.Println("Asc运行时间：", t4.Sub(t3))
		fmt.Println("Desc运行时间：", t5.Sub(t4))
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"asc":  resultAsc,
		"desc": resultDesc,
	})

}

/**
 * redis有序集合新增
 */
func (*Api) RedisZSet(c *gin.Context) {
	key, _ := c.GetQuery("key")
	score, _ := c.GetQuery("score")
	member, _ := c.GetQuery("member")

	sort, _ := strconv.Atoi(score)

	var redisClass redisLib.RedisInstanceClass
	redisClass.GetRedigoByName("default")
	result := redisClass.ZAdd(key, sort, member)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  result,
		"data": "",
	})
}

/**
 * redis有序集合删除
 */
func (*Api) RedisZRem(c *gin.Context) {
	key, _ := c.GetQuery("key")
	member, _ := c.GetQuery("member")

	var redisClass redisLib.RedisInstanceClass
	redisClass.GetRedigoByName("default")
	result := redisClass.ZRem(key, member)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  result,
		"data": "",
	})
}

/**
 * redis有序集合获取
 */
func (*Api) RedisZRange(c *gin.Context) {
	key, _ := c.GetQuery("key")
	start, _ := c.GetQuery("start")
	end, _ := c.GetQuery("end")

	s, _ := strconv.Atoi(start)
	e, _ := strconv.Atoi(end)
	var redisClass redisLib.RedisInstanceClass
	redisClass.GetRedigoByName("default")
	resultAsc := redisClass.ZRange(key, s, e)
	resultDesc := redisClass.ZRevrange(key, s, e)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "result",
		"data": map[string]interface{}{"asc": resultAsc, "desc": resultDesc},
	})
}

/**
 * es操作,put
 */
func (*Api) ElasticPut(c *gin.Context) {

	var (
		put       interface{}
		err       error
		ctx       = context.Background()
		monConn   *mongodbMod.MyGinAd
		mongoData []mongodbMod.MyGinData
	)

	if err = monConn.Mongodb().Find(bson.M{}).Skip(0).Limit(1000).All(&mongoData); err != nil {
		fmt.Println(err.Error())
	}

	esClient := elastic.Init()
	for _, data := range mongoData {
		if put, err = esClient.Index().Index("mongo").BodyJson(data).Do(ctx); err != nil {
			fmt.Println(err.Error())
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"data": put,
	})
}

/**
 * es操作,search
 */
func (*Api) ElasticSearch(c *gin.Context) {
	stimeQuery := c.DefaultQuery("stime", carbon.Now().DateString())
	etimeQuery := c.DefaultQuery("etime", carbon.Now().DateString())
	stime, _ := carbon.Parse(carbon.DefaultFormat, stimeQuery, "Asia/Shanghai")
	etime, _ := carbon.Parse(carbon.DefaultFormat, etimeQuery, "Asia/Shanghai")
	shour := stime.Local()
	ehour := etime.Local()
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "15"))

	fmt.Println(shour, ehour)

	var (
		ctx           = context.Background()
		testIndexName = "my_gin"
		item          []interface{}
		// after = []interface{}{"hour", "ad_id"}
	)

	offset := (page - 1) * limit

	esClient := elastic.Init()
	termQuery := elastic2.NewRangeQuery("hour").Gte(shour).Lte(ehour)
	searchResult, err := esClient.Search().
		TrackTotalHits(true). // 返回正确的总数
		// SearchAfter(after...).
		Index(testIndexName). // 搜索的索引名称
		Query(termQuery).     // 条件
		// SortBy(elastic2.NewFieldSort("hour").Desc(), elastic2.NewFieldSort("ad_id.keyword").Desc()).
		Sort("hour", false).          // 字符串排序加keyword
		Sort("ad_id.keyword", false). // 字符串排序加keyword
		From(offset).Size(limit).     // 分页
		Pretty(true).                 // pretty print request and response JSON
		Do(ctx)                       // 执行

	if err != nil {
		fmt.Println(err)
	}
	if searchResult.Hits == nil {
		fmt.Println("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.TotalHits() != 3 {
		fmt.Printf("expected SearchResult.TotalHits() = %d; got %d\n", 3, searchResult.TotalHits())
	}
	if len(searchResult.Hits.Hits) != 3 {
		fmt.Printf("expected len(SearchResult.Hits.Hits) = %d; got %d\n", 3, len(searchResult.Hits.Hits))
	}

	for _, hit := range searchResult.Hits.Hits {
		if hit.Index != testIndexName {
			fmt.Printf("expected SearchResult.Hits.Hit.Index = %q; got %q\n", testIndexName, hit.Index)
		}
		item = append(item, hit.Source)
	}

	total := searchResult.Hits.TotalHits.Value

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "",
		"total": total,
		"data":  item,
		// "searchResult": searchResult,
	})

}

/*
 * es操作，delete
 */
func (*Api) ElasticDelete(c *gin.Context) {
	var (
		ctx           = context.Background()
		testIndexName = "mongo"
		item          []mongodbMod.MyGinData
	)

	esClient := elastic.Init()

	searchResult, err := esClient.Search().
		Index(testIndexName). // search in index "twitter" er" field, ascending
		From(0).Size(2000).
		Pretty(true). // pretty print request and response JSON
		Do(ctx)       // execute
	if err != nil {
		fmt.Println(err)
	}
	if searchResult.Hits == nil {
		fmt.Println("expected SearchResult.Hits != nil; got nil")
	}
	if searchResult.TotalHits() != 3 {
		fmt.Printf("expected SearchResult.TotalHits() = %d; got %d\n", 3, searchResult.TotalHits())
	}
	if len(searchResult.Hits.Hits) != 3 {
		fmt.Printf("expected len(SearchResult.Hits.Hits) = %d; got %d\n", 3, len(searchResult.Hits.Hits))
	}

	for _, hit := range searchResult.Hits.Hits {
		if hit.Index != testIndexName {
			fmt.Printf("expected SearchResult.Hits.Hit.Index = %q; got %q\n", testIndexName, hit.Index)
		}
		var data mongodbMod.MyGinData
		json.Unmarshal(hit.Source, &data)
		item = append(item, data)
	}

	for k, _ := range item {
		res, err := esClient.Delete().
			Index(testIndexName).
			Id(strconv.Itoa(k + 1)).
			Do(context.Background())
		if err != nil {
			println(err.Error())
			return
		}
		fmt.Printf("delete result %s\n", res.Result)
	}

}

// GRPC
type GrpcApiServer struct {
}

// grpc服务端：redis zset 设置
func (s *GrpcApiServer) RedisZSet(ctx context.Context, req *SetReq) (resp *SetResp, err error) {
	key := req.Key
	score := req.Score
	member := req.Member
	scoreStr := strconv.FormatInt(score, 10)
	scoreInt, _ := strconv.Atoi(scoreStr)

	var redisClass redisLib.RedisInstanceClass
	redisClass.GetRedigoByName("default")

	b := new(SetResp)
	b.Data = redisClass.ZAdd(key, scoreInt, member)
	return b, nil
}

// grpc服务端：redis zset 获取
func (s *GrpcApiServer) RedisZRange(ctx context.Context, req *RangeReq) (resp *RangeResp, err error) {
	key := req.Key
	start := req.Start
	end := req.End

	startStr := strconv.FormatInt(start, 10)
	startInt, _ := strconv.Atoi(startStr)

	endStr := strconv.FormatInt(end, 10)
	endInt, _ := strconv.Atoi(endStr)

	var redisClass redisLib.RedisInstanceClass
	redisClass.GetRedigoByName("default")
	resultAsc := redisClass.ZRange(key, startInt, endInt)
	b := new(RangeResp)
	b.Data = resultAsc
	return b, nil
}

// grpc客户端：redis zset 设置
func (*Api) GrpcRedisZRange(c *gin.Context) {
	key, _ := c.GetQuery("key")
	start, _ := c.GetQuery("start")
	end, _ := c.GetQuery("end")

	conn, err := grpc.Dial("127.0.0.1:"+config.ParseYaml().Grpc_port, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("fail to dial: %v", err))
	}
	defer conn.Close()

	apiClient := NewTestApiClient(conn)

	data := new(RangeReq)

	data.Key = key
	data.Start, _ = strconv.ParseInt(start, 10, 64)
	data.End, _ = strconv.ParseInt(end, 10, 64)

	bi, _ := apiClient.RedisZRange(context.Background(), data)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"data": bi.Data,
	})
}

// grpc客户端：redis zset 获取
func (*Api) GrpcRedisZSet(c *gin.Context) {
	key, _ := c.GetQuery("key")
	score, _ := c.GetQuery("score")
	member, _ := c.GetQuery("member")

	conn, err := grpc.Dial("127.0.0.1:"+config.ParseYaml().Grpc_port, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("fail to dial: %v", err))
	}
	defer conn.Close()

	apiClient := NewTestApiClient(conn)

	data := new(SetReq)

	data.Key = key
	data.Score, _ = strconv.ParseInt(score, 10, 64)
	data.Member = member

	bi, _ := apiClient.RedisZSet(context.Background(), data)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"data": bi.Data,
	})
}
