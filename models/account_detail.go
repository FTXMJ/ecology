package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

//交易记录表
type AccountDetail struct {
	Id             int     `orm:"column(id);pk;auto"`
	UserId         string  `orm:column(user_id)`
	CurrentRevenue float64 `orm:column(current_revenue)` //本期收入
	CurrentOutlay  float64 `orm:column(current_outlay)`  //本期支出
	OpeningBalance float64 `orm:column(opening_balance)` //上期余额
	CurrentBalance float64 `orm:column(current_balance)` //本期余额
	CreateDate     string  `orm:column(create_date)`     //创建时间
	Comment        string  `orm:column(comment)`         //备注
	TxId           string  `orm:column(tx_id)`           //任务id
	Account        int     `orm:column(account)`         //生态仓库id
	CoinType       string  `orm:"column(coin_type)"`     // 币种信息
}

func (this *AccountDetail) TableName() string {
	return "account_detail"
}

func (this *AccountDetail) Insert() error {
	_, err := NewOrm().Insert(this)
	return err
}

func (this *AccountDetail) Update() (err error) {
	_, err = NewOrm().Update(this)
	return err
}

func FindLimitOneAndSaveAcc_d(o orm.Ormer, user_id, comment, tx_id string, money_out, money_in float64, account_id int) error {
	account_old := AccountDetail{}
	o.QueryTable("account_detail").
		Filter("user_id", user_id).
		Filter("account", account_id).
		OrderBy("-create_date").
		Limit(1).
		One(&account_old)
	if account_old.Id == 0 {
		account_old.CurrentBalance = 0
	}
	account_new := AccountDetail{
		UserId:         user_id,
		CurrentRevenue: money_in,
		CurrentOutlay:  money_out,
		OpeningBalance: account_old.CurrentBalance,
		CurrentBalance: account_old.CurrentBalance + money_in - money_out,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        comment,
		TxId:           tx_id,
		Account:        account_id,
	}
	if account_new.CurrentBalance < 0 {
		account_new.CurrentBalance = 0
	}
	_, err_acc := o.Insert(&account_new)
	if err_acc != nil {
		o.Rollback()
		return err_acc
	}

	_, err_txid := o.QueryTable("tx_id_list").Filter("tx_id", tx_id).Update(orm.Params{"state": "true"})
	if err_txid != nil {
		o.Rollback()
		return err_txid
	}

	_, err_up := o.QueryTable("account").Filter("id", account_id).Update(orm.Params{"balance": account_new.CurrentBalance})
	if err_up != nil {
		o.Rollback()
		return err_acc
	}
	return nil
}

func NewCreateAndSaveAcc_d(o orm.Ormer, user_id, comment, tx_id string, money_out, money_in float64, account_id int) error {
	o.Begin()

	account_new := AccountDetail{
		UserId:         user_id,
		CurrentRevenue: money_in,
		CurrentOutlay:  money_out,
		OpeningBalance: 0,
		CurrentBalance: money_in,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        comment,
		TxId:           tx_id,
		Account:        account_id,
	}

	_, err_acc := o.Insert(&account_new)
	if err_acc != nil {
		o.Rollback()
		return err_acc
	}

	_, err_txid := o.QueryTable("tx_id_list").Filter("tx_id", tx_id).Update(orm.Params{"state": "true"})
	if err_txid != nil {
		o.Rollback()
		return err_txid
	}

	_, err_up := o.QueryTable("account").Filter("id", account_id).Update(orm.Params{"balance": account_new.CurrentBalance})
	if err_up != nil {
		o.Rollback()
		return err_acc
	}
	return nil
}
