package controllers

import (
	"ecology/common"
	db "ecology/db"
	"ecology/models"

	"github.com/astaxie/beego"
)

type FirstController struct {
	beego.Controller
}

type Ping struct {
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

// @Tags 给定实时数据
// @Accept  json
// @Produce json
// @Success 200____给定实时数据 {object} models.RealTimePriceTest
// @router /ticker [GET]
func (this *Ping) ShowRealTimePrice() {
	var (
		data *common.ResponseData
		o    = db.NewEcologyOrm()
	)
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	value := []models.QuoteTicker{}
	o.Raw("select * from quote_ticker").QueryRows(&value)

	data = common.NewResponse(value)
	return
}
