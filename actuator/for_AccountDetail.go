package actuator

import (
	"ecology/models"

	"github.com/astaxie/beego/orm"

	"time"
)

// 生成充值表的借贷记录
func FindLimitOneAndSaveAcc_d(o orm.Ormer, user_id, comment, tx_id string, money_out, money_in float64, account_id int) error {
	account_old := models.AccountDetail{}
	o.Raw("select * from account_detail where user_id=? and account=? order by create_date desc,id desc limit 1", user_id, account_id).QueryRow(&account_old)
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
		return err_acc
	}

	account := models.Account{
		UserId: user_id,
	}
	o.Read(&account, "user_id")
	_, err_up := o.Raw("update account set balance=? where id=?", account_new.CurrentBalance, account_id).Exec()
	if err_up != nil {
		return err_acc
	}
	return nil
}

// 创建第一条充值表的借贷记录
func NewCreateAndSaveAcc_d(o orm.Ormer, user_id, comment, tx_id string, money_out, money_in float64, account_id int) error {
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
		return err_acc
	}

	_, err_up := o.Raw("update account set balance=? where id=?", account_new.CurrentBalance, account_id).Exec()
	if err_up != nil {
		return err_acc
	}
	return nil
}
