package main

import (
	"ecology/controllers"
	_ "ecology/routers"
	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

var c = cron.New()

func main() {
	//c.AddFunc("0 0 2 1/1 * ? ", controllers.Test{}.DailyDividendAndRelease) //定时更新首页数据
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
