package defaultExecution

/**
 *	rabbitmq接收服务,测试同一个队列两个接收者
 */
import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"my-gin/app/models/mongodb"
	"my-gin/libraries/log"
	mongodb2 "my-gin/libraries/mongodb"
	"my-gin/libraries/rabbitmq"
)

func MonitorAdHourMqTwo() {
	var data []mongodb.MyGinData
	var conn *mongodb.MyGin

	logger := log.InitLog("monitorAdHourMq")
	//消息接收

	ch := rabbitmq.Init("my_vhost")

	//创建交换器
	//err := ch.ExchangeDeclare("st", "fanout", true, true, false, false, nil)

	// 使用默认交换器
	// 指定队列！
	q, err := ch.QueueDeclare(
		"adHour", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	rabbitmq.FailOnError(err, "Failed to declare a queue")

	// Fair dispatch 预取，每个工作方每次拿一个消息，确认后才拿下一次，缓解压力
	err = ch.Qos(
		1, // prefetch count
		// 待解释
		0,     // prefetch size
		false, // global
	)
	rabbitmq.FailOnError(err, "Failed to set QoS")

	// 消费根据队列名
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack   设置为真自动确认消息
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	rabbitmq.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			_ = json.Unmarshal(d.Body, &data)
			fmt.Printf("队列消耗2：%+v\n", data)

			for _, val := range data {
				count, _ := conn.Mongodb().Find(bson.M{"ad_id": val.Ad_id}).Count()
				fmt.Printf("队列消耗2：%d 校验重复\n", count)
				if count == 0 {
					// param := Struct2Map(val)
					val.Id = mongodb2.CreateId()
					logger.Info(val.Id)
					err = conn.Mongodb().Insert(val)
				} else {
					fmt.Printf("队列消耗2：%s 广告id重复\n", val.Ad_id)
				}

			}

			// 确认消息被收到！！如果为真的，那么同在一个channel,在该消息之前未确认的消息都会确认，适合批量处理
			// 真时场景：每十条消息确认一次，类似
			if err == nil {
				d.Ack(true)
			} else {
				logger.Errorf("msg", err)
			}

		}

		//forever<-true
	}()

	<-forever
}
