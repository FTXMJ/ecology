package models

import "github.com/astaxie/beego/orm"

// 数据库链接实例
var db orm.Ormer

// 全网收入参数
var NetIncome float64

func NewOrm() orm.Ormer {
	db = orm.NewOrm()
	return db
}
