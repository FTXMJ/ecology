package models

import (
	"ecology/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dsn := beego.AppConfig.String("mysql::db")
	orm.RegisterModel(
		new(models.Account),
		new(models.AccountDetail),
		new(models.BlockedDetail),
		new(models.CalculationPower),
		new(models.Formula),
		new(models.TxIdList),
		new(models.SuperForceTable),
		new(models.ForceTable),
		new(models.DailyDividendTasks),
		new(models.SuperPeerTable),
		new(models.PeerHistory),
		new(models.GlobalOperations),
		new(models.MrsfStateTable),
		new(models.User))
	orm.Debug = true // 是否开启调试模式 调试模式下会打印出sql语句
	er := orm.RegisterDataBase("default", "mysql", dsn, 100, 200)
	fmt.Println(er)
}
