package actuator

import (
	"ecology/models"
	"github.com/astaxie/beego/orm"
	"time"
)

// 生成充值表的借贷记录
func FindLimitOneAndSaveAcc_d(o orm.Ormer, user_id, comment, tx_id string, money_out, money_in float64, account_id int) error {
	account_old := models.AccountDetail{}
	o.QueryTable("account_detail").
		Filter("user_id", user_id).
		Filter("account", account_id).
		OrderBy("-create_date").
		Limit(1).
		One(&account_old)
	if account_old.Id == 0 {
		account_old.CurrentBalance = 0
	}
	account_new := models.AccountDetail{
		UserId:         user_id,
		CurrentRevenue: money_in,
		CurrentOutlay:  money_out,
		OpeningBalance: account_old.CurrentBalance,
		CurrentBalance: account_old.CurrentBalance + money_in - money_out,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        comment,
		TxId:           tx_id,
		Account:        account_id,
		CoinType:       "USDD",
	}
	if account_new.CurrentBalance < 0 {
		account_new.CurrentBalance = 0
	}
	_, err_acc := o.Insert(&account_new)
	if err_acc != nil {
		o.Rollback()
		return err_acc
	}

	account := models.Account{
		UserId: user_id,
	}
	o.Read(&account, "user_id")
	_, err_up := o.QueryTable("account").Filter("id", account_id).Update(orm.Params{"balance": account_new.CurrentBalance})
	if err_up != nil {
		o.Rollback()
		return err_acc
	}
	return nil
}

// 创建第一条充值表的借贷记录
func NewCreateAndSaveAcc_d(o orm.Ormer, user_id, comment, tx_id string, money_out, money_in float64, account_id int) error {
	o.Begin()

	account_new := models.AccountDetail{
		UserId:         user_id,
		CurrentRevenue: money_in,
		CurrentOutlay:  money_out,
		OpeningBalance: 0,
		CurrentBalance: money_in,
		CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
		Comment:        comment,
		TxId:           tx_id,
		Account:        account_id,
		CoinType:       "USDD",
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
