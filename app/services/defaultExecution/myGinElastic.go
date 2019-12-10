package defaultExecution

import (
	"context"
	"fmt"
	"github.com/jianfengye/collection"
	"github.com/uniplaces/carbon"
	"gopkg.in/mgo.v2/bson"
	"my-gin/app/models/mongodb"
	"my-gin/libraries/elastic"
	"strconv"
	"time"
)

func MyGinElastic() {
	var (
		monConn  *mongodb.MyGinAd
		err      error
		page     = 427
		limit    = 10000
		ctx      = context.Background()
		esClient = elastic.Init()
	)

	for {
		offset := (page - 1) * limit
		var mysqlData []mongodb.MyGinAdData
		t1 := time.Now()
		if err = monConn.Mongodb().Find(bson.M{}).Skip(offset).Limit(limit).All(&mysqlData); err != nil {
			fmt.Println(err.Error())
		}
		t2 := time.Now()
		fmt.Println("mongodb运行时间：", t2.Sub(t1))

		t3 := time.Now()
		collection.NewObjCollection(mysqlData).Each(func(item interface{}, key int) {
			v := item.(mongodb.MyGinAdData)
			if _, err = esClient.Index().Index("my_gin").BodyJson(v).Do(ctx); err != nil {
				fmt.Println(err.Error())
			}
		})
		t4 := time.Now()
		fmt.Println("es运行时间：", t4.Sub(t3))

		fmt.Println("page:" + strconv.Itoa(offset) + ",total:" + strconv.Itoa(len(mysqlData)) + ",time:" + time.Now().Format(carbon.DefaultFormat))
		if len(mysqlData) < limit {
			break
		}
		page++
	}
}
