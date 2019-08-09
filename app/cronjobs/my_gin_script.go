package cronjobs

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"my-gin/app/models/mongodb"
	"my-gin/app/models/mysql"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

//var Chs = make([] chan int, 2) 使用通道阻塞主线程
var wg sync.WaitGroup //官方推荐阻塞主线程方法
var day_num = 10      //天数
var hour_num = 24     //小时数
func My_gin_script() {

	fmt.Println("开始执行")
	fmt.Println(runtime.NumGoroutine())
	t := time.Now()

	var monConn *mongodb.MyGin
	var adData []mongodb.MyGinData

	err := monConn.Mongodb().Find(bson.M{"campaign_id": bson.M{"$ne": ""}, "product_id": bson.M{"$gt": 0}, "advertiser_id": bson.M{"$gt": 0}, "division_id": bson.M{"$gt": 0}}).All(&adData)
	fmt.Println(adData)
	if err != nil {
		fmt.Println(err.Error())
	}

	wg.Add(len(adData))
	//
	//fmt.Println(adData)
	//return

	for _, data := range adData {

		//Chs[n] = make(chan int)
		go createData(&wg, data)
		//createData(&wg)

	}

	//for _, ch := range(Chs) {
	//	<-ch
	//}
	fmt.Println(runtime.NumGoroutine())

	wg.Wait()
	elapsed := time.Since(t)

	fmt.Println("app elapsed:", elapsed)
	fmt.Println("结束执行")

}

func createData(wg *sync.WaitGroup, adData mongodb.MyGinData) {
	var data mysql.My_gin
	var day string
	//定义某一个广告的基本数据

	data.Ad_id = adData.Id
	data.Campaign_id = adData.Campaign_id
	data.Product_id = adData.Product_id
	data.Advertiser_id = adData.Advertiser_id
	data.Division_id = adData.Division_id
	data.Status = 5
	data.Created_at = time.Now()
	data.Updated_at = time.Now()

	//记录多少天的数据
	for a := 1; a <= day_num; a++ {

		day = time.Date(2019, time.August, 30, 0, 0, 0, 0, time.UTC).AddDate(0, 0, a).Format("2006-01-02")

		//记录每天24小时的数据
		for h := 0; h < hour_num; h++ {
			randObj := rand.New(rand.NewSource(time.Now().UnixNano()))
			data.Request_count = randObj.Intn(10000)
			data.Cpm_count = randObj.Intn(1000) + 100
			data.Cpc_original_count = randObj.Intn(10)

			data.Hour = strings.Join([]string{day, strconv.Itoa(h)}, " ") + ":00:00"
			err := mysql.My_gin_obj().Create(&data).Error
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	wg.Done()
}
