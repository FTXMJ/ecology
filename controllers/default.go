package controllers

import (
	"ecology/common"
	"github.com/astaxie/beego"
)

type FirstController struct {
	beego.Controller
}

// @Tags 心跳检测
// @Accept  json
// @Produce json
// @Success 200
// @router /check [GET]
func (this *FirstController) Check() {
	var data *common.ResponseData
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	data = common.NewResponse(nil)
	return
}
