package main

import (
	"ecology/actuator"
	"ecology/conf"
	"ecology/controllers"
	"ecology/filter"
	"ecology/routers"
	_ "ecology/routers"
	"github.com/robfig/cron"
)

var c = cron.New()

func main() {
	//定时更新首页数据
	c.AddFunc(conf.ConfInfo.Schedules, controllers.DailyDividendAndRelease)
	c.AddFunc(conf.ConfInfo.Real_time_price, actuator.Second5s)

	c.Start()
	//http://localhost:8080/swagger/
	//bee run -gendoc=true -downdoc=true
	//bee pack -be GOOS=linux
	//	"github.com/shopspring/decv"
	filter.SetSignKey(conf.ConfInfo.Jwt)

	routers.Router()
}
