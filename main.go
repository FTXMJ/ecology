package main

import (
	"ecology/actuator"
	"ecology/conf"
	"ecology/controllers"
	"ecology/routers"
	_ "ecology/routers"
	"fmt"

	"github.com/robfig/cron"
)

var c = cron.New()

func main() {
	//定时更新首页数据
	c.AddFunc(conf.ConfInfo.Schedules, controllers.DailyDividendAndRelease)
	c.AddFunc(conf.ConfInfo.Real_time_price, actuator.Second5s)
	fmt.Println(conf.ConfInfo.Real_time_price)
	fmt.Println(conf.ConfInfo.Schedules)
	c.Start()
	//http://localhost:8080/swagger/
	//bee run -gendoc=true -downdoc=true
	//bee pack -be GOOS=linux
	//	"github.com/shopspring/decv"
	routers.Router()
}
