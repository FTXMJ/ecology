package main

import (
	"ecology/actuator"
	"ecology/controllers"
	"ecology/filter"
	_ "ecology/routers"
	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

var c = cron.New()

func main() {
	//定时更新首页数据
	c.AddFunc(beego.AppConfig.String("crontab::schedules"), controllers.DailyDividendAndRelease)
	c.AddFunc(beego.AppConfig.String("crontab::real_time_price"), actuator.Second5s)
	c.Start()
	//http://localhost:8080/swagger/
	//bee run -gendoc=true -downdoc=true
	//bee pack -be GOOS=linux
	//	"github.com/shopspring/decv"
	if beego.DEV == "dev" {
		beego.SetStaticPath("/api/v1/ecology/swagger", "swagger")
	}

	sk := beego.AppConfig.DefaultString("jwt::SignKey", "1233444")
	filter.SetSignKey(sk)
	beego.Run()
}
