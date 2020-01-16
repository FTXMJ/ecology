package filter

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	//拦截器验证登录
	beego.InsertFilter("/*", beego.BeforeExec, CheckLogin)
}

//验证登录
func CheckLogin(ctx *context.Context) {
	api := ctx.Request.URL.Path
	identity := ApiClassification(api)
	switch identity {
	case "admin":
		AdminFilter(ctx, api)
	case "user":
		UserFilter(ctx, api)
	default:
	}
}
