package controllers

import (
	"ecology1/common"
	"ecology1/models"
	"fmt"
	"github.com/astaxie/beego"
)

type FirstController struct {
	beego.Controller
}


// @Tags 测试空读
// @Accept  json
// @Produce json
// @Success 200___测试
// @router /test_read_all [Post]
func (this *FirstController) CreateUserAbout() {
	var data *common.ResponseData
	defer func() {
		this.Data["json"] = data
		this.ServeJSON()
	}()
	user := []models.User{}
	i,e := models.NewOrm().QueryTable("user").Filter("user_id","213").All(&user)
	fmt.Println("i := ",i)
	fmt.Println("e := ",e)
	 t := t{
		 A: i,
		 B: e,
	 }
	data = common.NewResponse(t)
}

type t struct {
	A int64
	B error
}