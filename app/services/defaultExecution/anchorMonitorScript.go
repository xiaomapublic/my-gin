package defaultExecution

import (
	"fmt"
	"math/rand"
	"my-gin/app/models/mysql"
	"runtime"
	"sync"
	"time"
)

const AGENTID = 167

var channelIds = []int{234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254, 255, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 285, 286, 287, 288, 289}

func AnchorMonitorScript() {
	fmt.Println("开始执行")
	fmt.Println(runtime.NumGoroutine())
	t := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(channelIds))

	for _, channelId := range channelIds {
		go createAnchorMonitor(&wg, channelId)
	}
	fmt.Println(runtime.NumGoroutine())
	wg.Wait()

	elapsed := time.Since(t)
	fmt.Println("执行时间:", elapsed)
	fmt.Println("结束执行")
}

func createAnchorMonitor(wg *sync.WaitGroup, channelId int) {
	dayNum := 365
	startTime := time.Date(2018, 9, 20, 0, 0, 0, 0, time.UTC)
	randObj := rand.New(rand.NewSource(time.Now().UnixNano()))

	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容
		}
		wg.Done()
	}()

	tx := mysql.AnchorMonitorObj().Begin()

	for i := 1; i <= dayNum; i++ {

		var anchorMonitor mysql.AnchorMonitor

		anchorMonitor.Date = startTime.AddDate(0, 0, i).Format("2006-01-02")
		anchorMonitor.Agent_id = AGENTID
		anchorMonitor.Agent_channel_id = channelId
		anchorMonitor.Cpm_count = randObj.Intn(1000)
		anchorMonitor.Income = (randObj.Float32() * 100) + 100
		anchorMonitor.Ecpm = anchorMonitor.Income / float32(anchorMonitor.Cpm_count) * 1000
		anchorMonitor.Status = 5
		anchorMonitor.Created_at = time.Now().Format("2006-01-02 15:04:05")
		anchorMonitor.Updated_at = time.Now().Format("2006-01-02 15:04:05")

		err := tx.Create(&anchorMonitor).Error
		if err != nil {
			tx.Rollback()
			panic(err.Error())
		}
	}
	tx.Commit()

}
