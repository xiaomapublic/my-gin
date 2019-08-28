package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	. "my-gin/libraries/config"
)

func Init(vhost string) *amqp.Channel {

	rabbitmq := DefaultConfig.GetStringMap("rabbitmq")
	conOne := rabbitmq[vhost].(map[string]interface{})
	addr := conOne["addr"].(string)
	user := conOne["user"].(string)
	pwd := conOne["pwd"].(string)
	return NewRabbitMq(vhost, addr, user, pwd)
}

func NewRabbitMq(vhost string, addr string, user string, pwd string) *amqp.Channel {
	//消息发布

	// 拨号，下面例子都一样
	conn, err := amqp.Dial("amqp://" + user + ":" + pwd + "@" + addr + "/" + vhost)
	FailOnError(err, "Failed to connect to RabbitMQ")

	// 这个是最重要的
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	return ch
}

func FailOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
