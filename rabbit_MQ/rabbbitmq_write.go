package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/streadway/amqp"
)

var channel *amqp.Channel
var conn *amqp.Connection

var notifyClose chan *amqp.Error

//init
func init() {
	if initChannel() {
		channel.NotifyClose(notifyClose)
	}
	// 断线自动重连
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				//log.Log.Printf("onNotifyChannelClosed: %+v\n", msg)
				initChannel()
			}
		}
	}()
}

func initChannel() bool {
	if channel != nil {
		return true
	}
	user_name := beego.AppConfig.String("rabbit_mq::user_name")
	user_passwd := beego.AppConfig.String("rabbit_mq::user_passwd")
	ip_port := beego.AppConfig.String("rabbit_mq::ip_port")
	conn, e := amqp.Dial("amqp://" + user_name + ":" + user_passwd + "@" + ip_port + "/")
	if e != nil {
		//log.Log.Println(e.Error())
		return false
	}
	channel, e = conn.Channel()
	if e != nil {
		//log.Log.Println(e.Error())
		return false
	}
	return true
}

func Publish(exchange, quName string, data interface{}) bool {
	fmt.Println("发送消息")
	bytes, _ := json.Marshal(data)
	_, err := channel.QueueDeclarePassive(quName, true, false, false, true, nil)
	if err != nil {
		// 队列不存在,声明队列
		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		_, err = channel.QueueDeclare(quName, true, false, false, true, nil)
		if err != nil {
			//log.Log.Printf("MQ注册队列失败:%s \n", err)
			return false
		}
	}
	err = channel.Publish(
		exchange,
		quName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         bytes,
		},
	)
	if err != nil {
		//log.Log.Println(err.Error())
		return false
	}
	return true
}
