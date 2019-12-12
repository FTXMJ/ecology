package controllers

import (
	"ecology1/common"
	"ecology1/models"
	"github.com/astaxie/beego"
)

type FirstController struct {
	beego.Controller
}


// @Tags 创建
// @Accept  json
// @Produce json
// @Success 200___测试
// @router /create_user [POST]
func (this *FirstController) CreateUserAbout() {

}

type t struct {
	A int64
	B error
}

// @Tags 心跳检测
// @Accept  json
// @Produce json
// @Success 200
// @router /check [GET]
func (this *FirstController)Check()  {
	var data *common.ResponseData
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	data = common.NewResponse(nil)
	return
}

// @Tags test
// @Accept  json
// @Produce json
// @Success 200
// @router /add_user [GET]
func (this *FirstController)AddUser()  {
	var data *common.ResponseData
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	user := models.User{
		Name:     "lxd",
		UserId:   "靓仔",
	}
	err := user.Insert()
	if err!=nil{
		data = common.NewErrorResponse(500)
		return
	}

	data = common.NewResponse(nil)
	return
}

// @Tags test
// @Accept  json
// @Produce json
// @Success 200
// @router /add_account [GET]
func (this *FirstController)AddAccount()  {
	var data *common.ResponseData
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	account := models.Account{
		UserId:        "lxd",
		Balance:       0,
		Currency:      "USDD",
		BockedBalance: 0,
		Level:         "",
	}
	err := account.Insert()
	if err!=nil{
		data = common.NewErrorResponse(500)
		return
	}

	data = common.NewResponse(nil)
	return
}

// @Tags test
// @Accept  json
// @Produce json
// @Success 200
// @router /add_formula [GET]
func (this *FirstController)AddFormula()  {
	var data *common.ResponseData
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()

	account := models.Formula{
		EcologyId:           53,
		Level:               "粉丝",
		LowHold:             0,
		HighHold:            0,
		ReturnMultiple:      0.004,
		HoldReturnRate:      0.09,
		RecommendReturnRate: 0.1,
		TeamReturnRate:      0.04,
	}
	err := account.Insert()
	if err!=nil{
		data = common.NewErrorResponse(500)
		return
	}

	data = common.NewResponse(nil)
	return
}