package main

import (
	"ecology/actuator"
	"ecology/controllers"
	"ecology/filter"
	"ecology/kafka"
	_ "ecology/routers"
	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

var c = cron.New()

func main() {
	//开启 kafka ,监听队列消息,并处理
	go kafka.AllTheTimeListen()

	//定时 每日释放
	c.AddFunc(beego.AppConfig.String("crontab::schedules"), controllers.DailyDividendAndRelease)

	//定时 行情更新
	c.AddFunc(beego.AppConfig.String("crontab::real_time_price"), actuator.Second5s)
	c.Start()

	if beego.DEV == "dev" {
		beego.SetStaticPath("/api/v1/ecology/swagger", "swagger")
	}

	sk := beego.AppConfig.DefaultString("jwt::SignKey", "1233444")
	filter.SetSignKey(sk)
	beego.Run()
}

//http://localhost:8080/swagger/
//bee run -gendoc=true -downdoc=true
//bee pack -be GOOS=linux
//	"github.com/shopspring/decv"
