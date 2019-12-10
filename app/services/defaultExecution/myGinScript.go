package defaultExecution

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/uniplaces/carbon"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"my-gin/app/models/mongodb"
	"my-gin/app/models/mysql"
	mysqlLib "my-gin/libraries/mysql"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// var Chs = make([] chan int, 2) 使用通道阻塞主线程

var day_num = 360 // 天数
var hour_num = 24 // 小时数
var monAdConn *mongodb.MyGinAd
var mysqlDB = mysqlLib.GetORMByName("my_gin")
var tiDB = mysqlLib.GetORMByName("tidb_gin")

func MyGinScript() {

	fmt.Println("开始执行")
	fmt.Println(runtime.NumGoroutine())
	t := time.Now()

	var monConn *mongodb.MyGin
	var adData []mongodb.MyGinData

	err := monConn.Mongodb().Find(bson.M{"campaign_id": bson.M{"$ne": ""}, "product_id": bson.M{"$gt": 0}, "advertiser_id": bson.M{"$gt": 0}, "division_id": bson.M{"$gt": 0}}).All(&adData)
	if err != nil {
		fmt.Println(err.Error())
	}
	createData(adData)
	fmt.Println(runtime.NumGoroutine())

	elapsed := time.Since(t)

	fmt.Println("app elapsed:", elapsed)
	fmt.Println("结束执行")

}

func createData(adDatas []mongodb.MyGinData) {
	var day string
	// 定义某一个广告的基本数据
	// 组装数据
	for a := 1; a <= day_num; a++ {
		day = time.Date(2019, time.June, 30, 0, 0, 0, 0, time.UTC).AddDate(0, 0, a).Format("2006-01-02")
		// 记录每天24小时的数据
		for h := 0; h < hour_num; h++ {
			var data mysql.MyGin
			insertData := make([]mysql.MyGin, 0, len(adDatas))
			dataObj, _ := carbon.Parse(carbon.DefaultFormat, strings.Join([]string{day, strconv.Itoa(h)}, " ")+":00:00", "Asia/Shanghai")
			data.Hour = dataObj.Local()
			randObj := rand.New(rand.NewSource(time.Now().UnixNano()))
			data.Status = 5
			data.Created_at = time.Now()
			data.Updated_at = time.Now()
			for _, adData := range adDatas {
				data.Ad_id = adData.Ad_id
				data.Campaign_id = adData.Campaign_id
				data.Product_id = adData.Product_id
				data.Advertiser_id = adData.Advertiser_id
				data.Division_id = adData.Division_id
				data.Request_count = randObj.Intn(9999999)
				data.Cpm_count = randObj.Intn(9999999)
				data.Cpc_original_count = randObj.Intn(9999999)

				insertData = append(insertData, data)
			}

			var err error
			var wg sync.WaitGroup // 官方推荐阻塞主线程方法
			wg.Add(3)
			go func(insertData []mysql.MyGin) {
				t1 := time.Now()
				if err = Creates(mysqlDB, insertData); err != nil {
					fmt.Println(err)
				}
				t2 := time.Now()
				fmt.Println("mysql运行时间：", t2.Sub(t1))
				wg.Done()

			}(insertData)

			go func(insertData []mysql.MyGin) {
				t1 := time.Now()
				if err = Creates(tiDB, insertData); err != nil {
					fmt.Println(err)
				}
				t2 := time.Now()
				fmt.Println("tidb运行时间：", t2.Sub(t1))
				wg.Done()
			}(insertData)

			go func(insertData []mysql.MyGin) {

				doc := make([]interface{}, 0, len(insertData))
				for _, v := range insertData {
					doc = append(doc, v)
				}
				t1 := time.Now()
				// 将切片用...打散输入
				if err = monAdConn.Mongodb().Insert(doc...); err != nil {
					fmt.Println(err)
				}
				t2 := time.Now()
				fmt.Println("mongodb运行时间：", t2.Sub(t1))
				wg.Done()
			}(insertData)

			wg.Wait()

		}
	}
}

func Creates(db *gorm.DB, data []mysql.MyGin) error {
	sql := "INSERT INTO `my_gin` (`hour`,`ad_id`,`campaign_id`,`product_id`,`advertiser_id`,`request_count`,`cpm_count`,`cpc_original_count`,`division_id`,`status`,`created_at`,`updated_at`) VALUES "
	// 循环data数组,组合sql语句
	for key, v := range data {
		if len(data)-1 == key {
			// 最后一条数据 以分号结尾
			sql += fmt.Sprintf("('%s','%s','%s',%d,%d,%d,%d,%d,%d,%d,'%s','%s');", v.Hour.Format(carbon.DefaultFormat), v.Ad_id, v.Campaign_id, v.Product_id, v.Advertiser_id, v.Request_count, v.Cpm_count, v.Cpc_original_count, v.Division_id, v.Status, v.Created_at.Format(carbon.DefaultFormat), v.Updated_at.Format(carbon.DefaultFormat))
		} else {
			sql += fmt.Sprintf("('%s','%s','%s',%d,%d,%d,%d,%d,%d,%d,'%s','%s'),", v.Hour.Format(carbon.DefaultFormat), v.Ad_id, v.Campaign_id, v.Product_id, v.Advertiser_id, v.Request_count, v.Cpm_count, v.Cpc_original_count, v.Division_id, v.Status, v.Created_at.Format(carbon.DefaultFormat), v.Updated_at.Format(carbon.DefaultFormat))
		}
	}
	return db.Exec(sql).Error
}
