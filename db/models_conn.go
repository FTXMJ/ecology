package models

import (
	"ecology/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// 数据库链接实例
var Db_ecology orm.Ormer
var Db_wallet orm.Ormer

// 全网收入参数
var NetIncome float64

func NewEcologyOrm() orm.Ormer {
	Db_ecology = orm.NewOrm()
	Db_ecology.Using("default")
	return Db_ecology
}

func NewWalletOrm() orm.Ormer {
	Db_wallet = orm.NewOrm()
	Db_wallet.Using("wallet")
	return Db_wallet
}

func init() {
	ds_ecology := beego.AppConfig.String("mysql::db_ecology")
	ds_wallet := beego.AppConfig.String("mysql::db_wallet")
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
		new(models.QuoteTickerHistory),
		new(models.QuoteTicker),
		new(models.WtQuote),
		new(models.DappTable),
		new(models.User))
	orm.Debug = true // 是否开启调试模式 调试模式下会打印出sql语句

	orm.RegisterDriver("mysql", orm.DRMySQL)
	if err := orm.RegisterDataBase("default", "mysql", ds_ecology, 100, 200); err != nil {
		beego.Emergency("Can't register db, err :", err)
	}
	Db_ecology = orm.NewOrm()
	Db_ecology.Using("default")

	if err := orm.RegisterDataBase("wallet", "mysql", ds_wallet, 60, 100); err != nil {
		beego.Emergency("Can't register db, err :", err)
	}
	Db_wallet = orm.NewOrm()
	Db_wallet.Using("wallet")
}
