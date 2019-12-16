package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var db orm.Ormer

func NewOrm() orm.Ormer {
	db = orm.NewOrm()
	return db
}

func init() {
	dsn := beego.AppConfig.String("mysql::db")
	//if dsn == "" {
	//	dsn = "root:Passwd123!@(101.132.79.34:3306)/ecology?charset=utf8&loc=Asia%2FShanghai&loc=Local"
	//}
	orm.RegisterModel(
		new(Account),
		new(AccountDetail),
		new(BlockedDetail),
		new(CalculationPower),
		new(Formula),
		new(TxIdList),
		new(SuperForceTable),
		new(ForceTable),
		new(DailyDividendTasks),
		new(SuperPeerTable),
		new(User))
	orm.Debug = true // 是否开启调试模式 调试模式下会打印出sql语句
	orm.RegisterDataBase("default", "mysql", dsn, 3000, 3000)
}
