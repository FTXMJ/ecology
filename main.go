package main

import (
	"ecology/controllers"
	"ecology/logs"
	_ "ecology/routers"
	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

var c = cron.New()

func main() {
	c.AddFunc(beego.AppConfig.String("crontab::schedules"), controllers.DailyDividendAndRelease) //定时更新首页数据
	//c.AddFunc("0/10 * * * * *", Fp)                                                              //定时更新首页数据
	c.Start()
	//http://localhost:8080/swagger/
	//bee run -gendoc=true -downdoc=true
	//bee pack -be GOOS=linux
	if beego.DEV == "dev" {
		beego.SetStaticPath("/api/v1/ecology/swagger", "swagger")
	}

	sk := beego.AppConfig.DefaultString("jwt::SignKey", "1233444")
	controllers.SetSignKey(sk)
	beego.Run()
}

func Fp() {
	logs.Log.Error("hello -_- world")
}
