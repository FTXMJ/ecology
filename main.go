package main

import (
	"ecology1/controllers"
	"github.com/astaxie/beego"
	_ "ecology1/routers"
)

func main() {
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