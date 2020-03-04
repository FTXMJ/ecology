package models

import (
	"ecology/conf"
	"ecology/models"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	// 数据库链接实例
	db_ecology *gorm.DB
	db_wallet  *gorm.DB

	// 全网收入参数
	NetIncome float64
)

func NewEcologyOrm() *gorm.DB {
	return db_ecology.New()
}

func NewWalletOrm() *gorm.DB {
	return db_wallet.New()
}

func init() {
	var err error
	//连接数据库
	if db_ecology, err = gorm.Open("mysql", conf.ConfInfo.MysqlEcology); err != nil {
		logrus.Errorf("mysql-ecology client link failure : %s", err)
	}
	db_ecology.SingularTable(true)
	db_ecology.DB().SetMaxIdleConns(100)
	db_ecology.DB().SetMaxOpenConns(200)
	db_ecology.LogMode(true)

	if db_wallet, err = gorm.Open("mysql", conf.ConfInfo.MysqlWallet); err != nil {
		logrus.Errorf("mysql-wallet client link failure : %s", err)
	}
	db_wallet.SingularTable(true)
	db_wallet.DB().SetMaxIdleConns(100)
	db_wallet.DB().SetMaxOpenConns(200)
	db_wallet.LogMode(true)

	db_ecology.AutoMigrate(
		&models.Account{},
		&models.AccountDetail{},
		&models.BlockedDetail{},
		&models.CalculationPower{},
		&models.Formula{},
		&models.TxIdList{},
		&models.SuperForceTable{},
		&models.ForceTable{},
		&models.DailyDividendTasks{},
		&models.SuperPeerTable{},
		&models.PeerHistory{},
		&models.GlobalOperations{},
		&models.MrsfStateTable{},
		&models.QuoteTickerHistory{},
		&models.QuoteTicker{},
		&models.DappTable{},
		&models.User{},
	)

	db_wallet.AutoMigrate(
		&models.WtQuote{},
	)
}
