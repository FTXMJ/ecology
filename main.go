package main

import (
	"ecology1/controllers"
	"github.com/astaxie/beego"
	_ "ecology1/routers"
	"github.com/robfig/cron"
)

var c = cron.New()

func main() {
	c.AddFunc("0 0 2 1/1 * ? ",controllers.WelfarePayment) //定时更新首页数据
	c.Start()
	//http://localhost:8080/swagger/
	//bee run -gendoc=true -downdoc=true
	//bee pack -be GOOS=linux
	if beego.DEV == "dev" {
		beego.SetStaticPath("/swagger", "swagger")
	}


	sk:= beego.AppConfig.DefaultString("jwt::SignKey","1233444")
	controllers.SetSignKey(sk)

	beego.Run(":2019")
}