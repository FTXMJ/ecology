package rabbitmq

import (
	"fmt"
)

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

var done chan bool

func Consumer(quName, cName string, callback func(msg []byte) bool) {
	fmt.Println("开始收取消息")
	_, err := channel.QueueDeclarePassive(quName, true, false, false, true, nil)
	if err != nil {
		//log.Log.Info("user队列不存在")
		// 队列不存在,声明队列
		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		_, err = channel.QueueDeclare(quName, true, false, false, true, nil)
		if err != nil {
			//log.Log.Printf("MQ注册队列失败:%s \n", err)
		}
	}
	channel.Qos(1, 0, true)
	consume, err := channel.Consume(
		quName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		//log.Log.Println(err.Error())
		return
	}

	done = make(chan bool)
	fmt.Println("开始收取消息")
	go func() {
		// 循环读取channel的数据
		for d := range consume {
			processErr := callback(d.Body)
			if !processErr {
				// TODO: 将任务写入错误队列，待后续处理
				//log.Log.Info("消费失败")
			} else {
				d.Ack(true)
			}
		}
	}()

	// 接收done的信号, 没有信息过来则会一直阻塞，避免该函数退出
	<-done
	// 关闭通道
	channel.Close()

}

// StopConsume : 停止监听队列
func StopConsume() {
	done <- true
}
