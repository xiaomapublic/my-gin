package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	. "my-gin/app/libraries/config"
)

var RabbitSession map[string]*amqp.Channel

func Init() {

	RabbitSession = make(map[string]*amqp.Channel, len(DefaultConfig.GetStringMap("rabbitmq")))

	for key, c := range DefaultConfig.GetStringMap("rabbitmq") {
		conOne := c.(map[string]interface{})
		vhost := key
		addr := conOne["addr"].(string)
		user := conOne["user"].(string)
		pwd := conOne["pwd"].(string)
		RabbitSession[vhost] = NewRabbitMq(vhost, addr, user, pwd)
	}

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
